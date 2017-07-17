// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nodescaler

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	core "k8s.io/client-go/testing"
)

type NodePoolMock struct {
	size int64
}

// IncreaseToSize
func (npm *NodePoolMock) IncreaseToSize(size int64) error {
	if size <= npm.size {
		log.Printf("[Test][Mock:IncreaseToSize] Ignoring resize to %v, as size is already %v", size, npm.size)
		return nil
	}
	log.Printf("[Test][Mock:IncreaseToSize] Resising to: %v", size)
	npm.size = size
	return nil
}

func TestScaleUpNodes(t *testing.T) {
	t.Parallel()

	nodes := newNodeListFixture(nlConfig{count: 2, cpu: []string{"2.0", "2.0"}})
	assertAllUnscheduled(t, nodes)

	cs := &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)

	s, err := NewServer("", "app=game-server", "0.5", 5, time.Second)
	assert.Nil(t, err)
	s.cs = cs

	expected := int64(len(nodes.Items))
	mock := &NodePoolMock{size: expected}
	s.nodePool = mock
	err = s.scaleNodes()
	assert.Nil(t, err)
	assert.Equal(t, expected, mock.size)
	assertAllUnscheduled(t, nodes)

	pl1 := newPodListFixture([]string{"0.5", "0.3"})
	cs.AddReactor("list", "pods", func(a core.Action) (bool, runtime.Object, error) {
		if a.(core.ListAction).GetListRestrictions().Fields.String() == "spec.nodeName=node0" {
			return true, pl1, nil
		}
		return false, nil, nil
	})

	s.nodePool = mock
	err = s.scaleNodes()
	assert.Nil(t, err)
	assert.Equal(t, expected, mock.size)
	assertAllUnscheduled(t, nodes)

	pl2 := newPodListFixture([]string{"1.8"})
	cs.AddReactor("list", "pods", func(a core.Action) (bool, runtime.Object, error) {
		if a.(core.ListAction).GetListRestrictions().Fields.String() == "spec.nodeName=node1" {
			return true, pl2, nil
		}
		return false, nil, nil
	})

	s.nodePool = mock
	err = s.scaleNodes()
	assert.Nil(t, err)
	assert.Equal(t, int64(3), mock.size)
	assertAllUnscheduled(t, nodes)
}

func TestScaleUpCordonedNodesNoPods(t *testing.T) {
	t.Parallel()

	nodes := newNodeListFixture(nlConfig{count: 2, cpu: []string{"2.0", "2.0"}})
	nodes.Items[0].Spec.Unschedulable = true
	// gate it, just to be sure
	assert.True(t, nodes.Items[0].Spec.Unschedulable)
	assert.False(t, nodes.Items[1].Spec.Unschedulable)

	cs := &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)

	// use a size of 0, as a way to ensure this doesn't get called.
	mock := &NodePoolMock{size: 0}
	s, err := NewServer("", "app=game-server", "0.5", 5, time.Second)
	assert.Nil(t, err)
	s.cs = cs
	s.nodePool = mock

	err = s.scaleNodes()
	assert.Nil(t, err)
	assertAllUnscheduled(t, nodes)
	assert.Equal(t, int64(0), mock.size)

	// same test, but with three nodes
	nodes = newNodeListFixture(nlConfig{count: 3, cpu: []string{"2.0", "2.0", "2.0"}})
	nodes.Items[0].Spec.Unschedulable = true
	nodes.Items[1].Spec.Unschedulable = true

	mock = &NodePoolMock{size: 0}
	s, err = NewServer("", "app=game-server", "0.5", 5, time.Second)
	assert.Nil(t, err)
	cs = &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)
	s.cs = cs
	s.nodePool = mock

	err = s.scaleNodes()
	assert.Nil(t, err)
	nl, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(nl.availableNodes()))
	assert.Equal(t, 1, len(nl.cordonedNodes()))
	assert.Equal(t, int64(0), mock.size)
}

func TestScaleUpCordonedNodesWithPods(t *testing.T) {
	// we fake pods, by manipulating the allocatable value
	nodes := newNodeListFixture(nlConfig{count: 3, cpu: []string{"2.0", "2.0", "2.0"}})
	nodes.Items[0].Spec.Unschedulable = true
	nodes.Items[1].Spec.Unschedulable = true
	pods := newPodListFixture([]string{"0.5"})

	cs := &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)
	defaultListPodReactor(cs, map[string]*v1.PodList{"node1": pods})

	// use a size of 0, as a way to ensure this doesn't get called.
	mock := &NodePoolMock{size: 0}
	s, err := NewServer("", "app=game-server", "0.5", 5, time.Second)
	assert.Nil(t, err)
	s.cs = cs
	s.nodePool = mock

	err = s.scaleNodes()
	assert.Nil(t, err)
	assert.False(t, nodes.Items[2].Spec.Unschedulable, "Node2 should not be unscheduled")
	assert.False(t, nodes.Items[1].Spec.Unschedulable, "Node1 should not be unscheduled")
	assert.True(t, nodes.Items[0].Spec.Unschedulable, "Node0 should be unscheduled")
	assert.Equal(t, int64(0), mock.size)
}

func TestScaleDownCordonTwoNodes(t *testing.T) {
	nodes := newNodeListFixture(nlConfig{count: 2, cpu: []string{"5.0", "5.0"}})

	cs := &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)

	s, err := NewServer("", "app=game-server", "0.5", 5, time.Second)
	assert.Nil(t, err)
	s.cs = cs

	// Make sure at least one of them is unscheduled
	err = s.scaleNodes()
	assert.Nil(t, err)

	nl, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nl.availableNodes()))
	assert.Equal(t, 1, len(nl.cordonedNodes()))

	// reset nodes, and up their capacity to 6.0
	for i, n := range nodes.Items {
		n.Spec.Unschedulable = false
		n.Status.Capacity.Cpu().Set(6)
		nodes.Items[i] = n
	}

	pl1 := newPodListFixture([]string{"0.5"})
	defaultListPodReactor(cs, map[string]*v1.PodList{"node0": {}, "node1": pl1})

	log.Printf("[Test][%v] Scaling down after resetting capacity to 6, and adding a single pod to node1.", t.Name())
	err = s.scaleNodes()
	assert.Nil(t, err)
	nl, err = s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nl.availableNodes()))
	assert.Equal(t, "node0", nl.nodes.Items[0].Name)
}

func TestScaleDownCordonThreeNodes(t *testing.T) {
	nodes := newNodeListFixture(nlConfig{count: 3, cpu: []string{"5.0", "5.0", "5.0"}})

	cs := &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)

	s, err := NewServer("", "app=game-server", "0.5", 5, time.Second)
	assert.Nil(t, err)
	s.cs = cs

	// gate - make sure everything is correct
	nl, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(nl.availableNodes()))
	assert.Equal(t, 0, len(nl.cordonedNodes()))

	err = s.scaleNodes()
	assert.Nil(t, err)

	// should only be 1 left, and now 2 should be cordoned
	nl, err = s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nl.availableNodes()))
	assert.Equal(t, 2, len(nl.cordonedNodes()))

	err = s.scaleNodes()
	assert.Nil(t, err)

	// make sure everything stays the same when running scaling again
	nl2, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, nl.availableNodes(), nl2.availableNodes())
	assert.Equal(t, nl.cordonedNodes(), nl2.cordonedNodes())
}
