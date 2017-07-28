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

	"github.com/jonboulle/clockwork"
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

// DeleteNode mock
func (npm *NodePoolMock) DeleteNodes(nodes []v1.Node) error {
	npm.size -= int64(len(nodes))
	return nil
}

func TestEnsureMinimumNumberOfNodes(t *testing.T) {
	t.Parallel()

	nodes := newNodeListFixture(nlConfig{count: 1, cpu: []string{"2.0"}})
	cs := &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)

	s, err := NewServer("", "app=game-server", "0.5", ServerMinNodeNumber(5))
	assert.Nil(t, err)
	s.cs = cs
	mock := &NodePoolMock{size: 1}
	s.nodePool = mock
	assert.EqualValues(t, 1, mock.size)

	err = s.scaleNodes()
	assert.Nil(t, err)
	assert.EqualValues(t, 5, mock.size)
}

func TestScaleUpNodes(t *testing.T) {
	t.Parallel()

	nodes := newNodeListFixture(nlConfig{count: 2, cpu: []string{"2.0", "2.0"}})
	assertAllUnscheduled(t, nodes)

	cs := &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)

	s, err := NewServer("", "app=game-server", "0.5", ServerBufferCount(5))
	assert.Nil(t, err)
	s.cs = cs

	expected := int64(len(nodes.Items))
	mock := &NodePoolMock{size: expected}
	s.nodePool = mock
	err = s.scaleNodes()
	assert.Nil(t, err)
	assert.Equal(t, expected, mock.size)
	assertAllUnscheduled(t, nodes)

	pl1 := newPodListFixture([]string{"0.5", "0.3"}, false)
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

	pl2 := newPodListFixture([]string{"1.8"}, false)
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

func TestScaleUpNodesToMax(t *testing.T) {
	t.Parallel()

	nodes := newNodeListFixture(nlConfig{count: 1, cpu: []string{"2.0"}})
	assertAllUnscheduled(t, nodes)

	cs := &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)

	s, err := NewServer("", "app=game-server", "0.5", ServerBufferCount(500), ServerMaxNodeNumber(5))
	assert.Nil(t, err)
	s.cs = cs

	mock := &NodePoolMock{size: 1}
	s.nodePool = mock
	err = s.scaleNodes()
	assert.Nil(t, err)
	assert.EqualValues(t, 5, mock.size)
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
	s, err := NewServer("", "app=game-server", "0.5", ServerBufferCount(5))
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
	s, err = NewServer("", "app=game-server", "0.5", ServerBufferCount(5))
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
	t.Parallel()
	// we fake pods, by manipulating the allocatable value
	nodes := newNodeListFixture(nlConfig{count: 3, cpu: []string{"2.0", "2.0", "2.0"}})
	nodes.Items[0].Spec.Unschedulable = true
	nodes.Items[1].Spec.Unschedulable = true
	pods := newPodListFixture([]string{"0.5"}, false)

	cs := &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)
	defaultListPodReactor(cs, map[string]*v1.PodList{"node1": pods})

	// use a size of 0, as a way to ensure this doesn't get called.
	mock := &NodePoolMock{size: 0}
	s, err := NewServer("", "app=game-server", "0.5", ServerBufferCount(5))
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
	t.Parallel()
	nodes := newNodeListFixture(nlConfig{count: 2, cpu: []string{"5.0", "5.0"}})

	cs := &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)
	mock := &NodePoolMock{size: 0}

	s, err := NewServer("", "app=game-server", "0.5", ServerBufferCount(5))
	assert.Nil(t, err)
	s.cs = cs
	s.nodePool = mock

	// Make sure at least one of them is unscheduled
	err = s.scaleNodes()
	assert.Nil(t, err)

	nl, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nl.availableNodes()))
	assert.Equal(t, 1, len(nl.cordonedNodes()))
	assert.EqualValues(t, 0, mock.size)

	// reset nodes, and up their capacity to 6.0
	for i, n := range nodes.Items {
		n.Spec.Unschedulable = false
		n.Status.Capacity.Cpu().Set(6)
		nodes.Items[i] = n
	}

	pl1 := newPodListFixture([]string{"0.5"}, false)
	defaultListPodReactor(cs, map[string]*v1.PodList{"node0": {}, "node1": pl1})

	log.Printf("[Test][%v] Scaling down after resetting capacity to 6, and adding a single pod to node1.", t.Name())
	err = s.scaleNodes()
	assert.Nil(t, err)
	nl, err = s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nl.availableNodes()))
	assert.Equal(t, "node0", nl.nodes.Items[0].Name)
	assert.EqualValues(t, 0, mock.size)
}

func TestScaleDownCordonThreeNodes(t *testing.T) {
	t.Parallel()
	nodes := newNodeListFixture(nlConfig{count: 3, cpu: []string{"5.0", "5.0", "5.0"}})

	cs := &fake.Clientset{}
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)
	mock := &NodePoolMock{size: 0}

	s, err := NewServer("", "app=game-server", "0.5", ServerBufferCount(5))
	assert.Nil(t, err)
	s.cs = cs
	s.nodePool = mock

	// gate - make sure everything is correct
	nl, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(nl.availableNodes()))
	assert.Equal(t, 0, len(nl.cordonedNodes()))

	err = s.scaleNodes()
	assert.Nil(t, err)
	assert.EqualValues(t, 0, mock.size)

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
	assert.EqualValues(t, 0, mock.size)
}

func TestRemoveNodeAfterCordonTwoNodesNoPods(t *testing.T) {
	nodes := newNodeListFixture(nlConfig{count: 2, cpu: []string{"5.0", "5.0"}})
	mock := &NodePoolMock{size: 2}
	cs := &fake.Clientset{}
	s, err := NewServer("", "app=game-server", "0.5", ServerBufferCount(5), ServerShutdown(time.Minute))
	assert.Nil(t, err)
	s.cs = cs
	fc := clockwork.NewFakeClock()
	s.clock = fc
	s.nodePool = mock
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)
	defaultDeleteNodeReactor(cs, nodes)

	// Make sure at least one of them is unscheduled
	err = s.scaleNodes()
	assert.Nil(t, err)

	nl, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nl.availableNodes()))
	assert.Equal(t, 1, len(nl.cordonedNodes()))
	assert.EqualValues(t, 2, mock.size)

	// forward 30 seconds, given the same set of pods, all should stay the same
	fc.Advance(30 * time.Second)
	err = s.scaleNodes()
	assert.Nil(t, err)

	// make sure everything stays the same when running scaling again
	nl2, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, nl.availableNodes(), nl2.availableNodes())
	assert.Equal(t, nl.cordonedNodes(), nl2.cordonedNodes())
	assert.EqualValues(t, 2, mock.size)

	// forward 31 seconds, then we should drop a node
	fc.Advance(31 * time.Second)
	err = s.scaleNodes()
	assert.Nil(t, err)
	nl, err = s.newNodeList()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(nl.nodes.Items), "There should be only one node left, the cordoned one should be deleted")
	assert.Equal(t, nl.availableNodes(), nl2.availableNodes(), "The cordoned node should be the one deleted.")
	assert.EqualValues(t, 1, mock.size, "A node should also be deleted")

	// forward 30 seconds, given the same set of pods, all should stay the same
	fc.Advance(30 * time.Second)
	err = s.scaleNodes()
	assert.Nil(t, err)

	// make sure everything stays the same when running scaling again
	nl2, err = s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, nl.availableNodes(), nl2.availableNodes())
	assert.Equal(t, nl.cordonedNodes(), nl2.cordonedNodes())
	assert.EqualValues(t, 1, mock.size)
}

func TestRemoveNodeAfterCordonTwoNodesWithNonGamePods(t *testing.T) {
	nodes := newNodeListFixture(nlConfig{count: 2, cpu: []string{"5.0", "5.0"}})
	mock := &NodePoolMock{size: 2}
	cs := &fake.Clientset{}
	s, err := NewServer("", "app=game-server", "0.5", ServerBufferCount(5), ServerShutdown(time.Minute))
	assert.Nil(t, err)
	s.cs = cs
	fc := clockwork.NewFakeClock()
	s.clock = fc
	s.nodePool = mock
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)
	defaultDeleteNodeReactor(cs, nodes)
	// add pods to both nodes, but they are non-game pods - everything should run the same.
	pl1 := newPodListFixture([]string{"0.5"}, false)
	defaultListPodReactor(cs, map[string]*v1.PodList{"node0": pl1, "node1": pl1})

	// Make sure at least one of them is unscheduled
	err = s.scaleNodes()
	assert.Nil(t, err)

	nl, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nl.availableNodes()))
	assert.Equal(t, 1, len(nl.cordonedNodes()))
	assert.EqualValues(t, 2, mock.size)

	// forward 30 seconds, given the same set of pods, all should stay the same
	fc.Advance(30 * time.Second)
	err = s.scaleNodes()
	assert.Nil(t, err)

	// make sure everything stays the same when running scaling again
	nl2, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, nl.availableNodes(), nl2.availableNodes())
	assert.Equal(t, nl.cordonedNodes(), nl2.cordonedNodes())
	assert.EqualValues(t, 2, mock.size)

	// forward 31 seconds, then we should drop a node
	fc.Advance(31 * time.Second)
	err = s.scaleNodes()
	assert.Nil(t, err)
	nl, err = s.newNodeList()
	assert.Nil(t, err)

	assert.Equal(t, 1, len(nl.nodes.Items), "There should be only one node left, the cordoned one should be deleted")
	assert.Equal(t, nl.availableNodes(), nl2.availableNodes(), "The cordoned node should be the one deleted.")
	assert.EqualValues(t, 1, mock.size, "A node should also be deleted")

	// forward 30 seconds, given the same set of pods, all should stay the same
	fc.Advance(30 * time.Second)
	err = s.scaleNodes()
	assert.Nil(t, err)

	// make sure everything stays the same when running scaling again
	nl2, err = s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, nl.availableNodes(), nl2.availableNodes())
	assert.Equal(t, nl.cordonedNodes(), nl2.cordonedNodes())
	assert.EqualValues(t, 1, mock.size)
}

func TestDoNotRemoveNodeAfterCordonTwoNodsWithPod(t *testing.T) {
	nodes := newNodeListFixture(nlConfig{count: 2, cpu: []string{"5.0", "5.0"}})
	mock := &NodePoolMock{size: 2}
	cs := &fake.Clientset{}
	s, err := NewServer("", "app=game-server", "0.5", ServerBufferCount(5), ServerShutdown(time.Minute))
	assert.Nil(t, err)
	s.cs = cs
	fc := clockwork.NewFakeClock()
	s.clock = fc
	s.nodePool = mock

	p0 := newPodListFixture([]string{"0.5"}, true)
	p1 := newPodListFixture([]string{"0.5"}, true)

	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)
	defaultDeleteNodeReactor(cs, nodes)
	defaultListPodReactor(cs, map[string]*v1.PodList{"node0": p0, "node1": p1})

	// Make sure at least one of them is unscheduled
	err = s.scaleNodes()
	assert.Nil(t, err)

	nl, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nl.availableNodes()))
	assert.Equal(t, 1, len(nl.cordonedNodes()))
	assert.EqualValues(t, 2, mock.size)

	// forward 30 seconds, given the same set of pods, all should stay the same
	fc.Advance(30 * time.Second)
	err = s.scaleNodes()
	assert.Nil(t, err)

	// make sure everything stays the same when running scaling again
	nl2, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, nl.availableNodes(), nl2.availableNodes())
	assert.Equal(t, nl.cordonedNodes(), nl2.cordonedNodes())
	assert.EqualValues(t, 2, mock.size)

	// forward 31 seconds, all should also stay the same
	fc.Advance(31 * time.Second)
	err = s.scaleNodes()
	nl2, err = s.newNodeList()
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.Equal(t, nl.availableNodes(), nl2.availableNodes())
	assert.Equal(t, nl.cordonedNodes(), nl2.cordonedNodes())
	assert.EqualValues(t, 2, mock.size)
}

func TestRemoveNodeWithMinimumNodeCount(t *testing.T) {
	nodes := newNodeListFixture(nlConfig{count: 5, cpu: []string{"5.0", "5.0", "5.0", "5.0", "5.0"}})
	mock := &NodePoolMock{size: 5}
	cs := &fake.Clientset{}
	s, err := NewServer("", "app=game-server", "0.5", ServerBufferCount(5), ServerShutdown(time.Minute), ServerMinNodeNumber(3))
	assert.Nil(t, err)
	s.cs = cs
	fc := clockwork.NewFakeClock()
	s.clock = fc
	s.nodePool = mock
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)
	defaultDeleteNodeReactor(cs, nodes)

	// Make sure 4 of them is unscheduled
	log.Print("[Debug][Test] Scaling nodes. Should cordon 4")
	err = s.scaleNodes()
	assert.Nil(t, err)

	nl1, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nl1.availableNodes()))
	assert.Equal(t, 4, len(nl1.cordonedNodes()))
	assert.EqualValues(t, 5, mock.size)

	// forward 61 seconds, then we should several nodes
	fc.Advance(61 * time.Second)
	log.Print("[Debug][Test] Forward 61 seconds, scaling. Should delete 2 nodes.")
	err = s.scaleNodes()
	assert.Nil(t, err)
	nl2, err := s.newNodeList()
	assert.Nil(t, err)

	assert.Equal(t, 3, len(nl2.nodes.Items), "There should be three left, as that is the in number")
	assert.EqualValues(t, 3, mock.size, "Two nodes should also be deleted")

	// forward 30 seconds, given the same set of pods, all should stay the same
	fc.Advance(30 * time.Second)
	log.Print("[Debug][Test] Going forward 30 seconds, and scaling nodes. All should be the same")
	err = s.scaleNodes()
	assert.Nil(t, err)

	// make sure everything stays the same when running scaling again
	nl1, err = s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, nl1.availableNodes(), nl2.availableNodes())
	assert.Equal(t, nl1.cordonedNodes(), nl2.cordonedNodes())
	assert.EqualValues(t, 3, mock.size)
}

func TestFilterGameSessionPods(t *testing.T) {
	// none of these have sessions=game as a label
	pods := newPodListFixture([]string{"0.5", "0.3", "0.2"}, false)
	result := filterGameSessionPods(pods.Items)
	assert.Equal(t, 0, len(result))

	labels := map[string]string{"sessions": "game"}
	pods.Items[0].ObjectMeta.Labels = labels
	pods.Items[1].ObjectMeta.Labels = labels

	result = filterGameSessionPods(pods.Items)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, pods.Items[0:2], result)
}
