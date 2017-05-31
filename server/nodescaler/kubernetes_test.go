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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
)

func TestSumCPUResourceRequests(t *testing.T) {
	t.Parallel()
	nl := newNodeListFixture(nlConfig{count: 2, cpu: []string{"4.0", "4.0"}})
	pods := newPodListFixture([]string{"0.5", "0.3"})
	fnl := nodeList{nodes: nl, pods: map[string]*v1.PodList{"node0": pods, "node1": {}}, cpuRequest: 500}

	expected := resource.MustParse("0.8")
	assert.Equal(t, expected.MilliValue(), fnl.sumResourceLimit(nl.Items[0], v1.ResourceCPU))
	expected = resource.MustParse("0")
	assert.Equal(t, expected.MilliValue(), fnl.sumResourceLimit(nl.Items[1], v1.ResourceCPU))
}

func TestNewNodeList(t *testing.T) {
	t.Parallel()
	nodes := newNodeListFixture(nlConfig{count: 2, cpu: []string{"4.0", "4.0"}})
	log.Printf("Nodes: %#v", nodes)

	pl1 := newPodListFixture([]string{"0.5", "0.3"})
	pl2 := newPodListFixture([]string{"1.2"})

	cs := &fake.Clientset{}
	s, err := NewServer("", "app=game-server", "0.5", 5, time.Second)
	assert.Nil(t, err)
	s.cs = cs
	defaultListNodeReactor(cs, nodes)
	defaultListPodReactor(cs, map[string]*v1.PodList{"node0": pl1, "node1": pl2})

	nodeList, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, nodes, nodeList.nodes)
	assert.Equal(t, pl1, nodeList.nodePods(nodes.Items[0]))
	assert.Equal(t, pl2, nodeList.nodePods(nodes.Items[1]))
}

func TestNodeReady(t *testing.T) {
	t.Parallel()
	n := v1.Node{ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Status: v1.NodeStatus{}}

	assert.False(t, nodeReady(n))

	n = v1.Node{ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Status: v1.NodeStatus{Conditions: []v1.NodeCondition{{Type: v1.NodeReady, Status: v1.ConditionFalse}}}}
	assert.False(t, nodeReady(n))

	n = v1.Node{ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Status: v1.NodeStatus{Conditions: readyNodeCondition}}
	assert.True(t, nodeReady(n))
}

func TestAvailableNodes(t *testing.T) {
	t.Parallel()
	nodes := newNodeListFixture(nlConfig{count: 4, cpu: []string{"1.0", "2.0", "3.0", "2.0"}})
	nodes.Items[2].Status.Conditions = []v1.NodeCondition{{Type: v1.NodeReady, Status: v1.ConditionFalse}}
	nodes.Items[3].Spec.Unschedulable = true

	nl := nodeList{nodes: nodes}
	expected := []v1.Node{nodes.Items[0], nodes.Items[1]}
	an := nl.availableNodes()

	assert.Equal(t, expected, an)
}

func TestCordonedNodes(t *testing.T) {
	t.Parallel()
	nodes := newNodeListFixture(nlConfig{count: 4, cpu: []string{"1.0", "2.0", "3.0", "2.0"}})
	nodes.Items[3].Spec.Unschedulable = true

	nl := nodeList{nodes: nodes}
	expected := []v1.Node{nodes.Items[3]}
	an := nl.cordonedNodes()

	assert.Equal(t, expected, an)
}

func TestCpuRequestsAvailable(t *testing.T) {
	t.Parallel()
	nodes := newNodeListFixture(nlConfig{count: 2, cpu: []string{"2.0", "2.0"}})

	cs := &fake.Clientset{}
	s, err := NewServer("", "app=game-server", "0.5", 5, time.Second)
	assert.Nil(t, err)
	s.cs = cs
	defaultListNodeReactor(cs, nodes)

	nl, err := s.newNodeList()
	assert.Nil(t, err)
	count := nl.cpuRequestsAvailable()
	assert.Equal(t, int64(8), count)

	pl1 := newPodListFixture([]string{"0.5", "0.3"})
	pl2 := newPodListFixture([]string{"1.8"})
	defaultListPodReactor(cs, map[string]*v1.PodList{"node0": pl1, "node1": pl2})
	nl, err = s.newNodeList()
	assert.Nil(t, err)

	count = nl.cpuRequestsAvailable()
	assert.Equal(t, int64(2), count)
}

func TestNewGameWatcher(t *testing.T) {
	t.Parallel()

	mw := watch.NewFake()
	gw := &gameWatcher{events: make(chan bool), watcher: mw}
	gw.start()
	go func() {
		defer gw.stop()
		mw.Action(watch.Added, nil)
		mw.Action(watch.Deleted, nil)
		mw.Action(watch.Error, nil)
		mw.Action(watch.Modified, nil)
		mw.Action(watch.Deleted, nil)
	}()

	i := 0
	for range gw.events {
		i++
	}
	assert.Equal(t, 3, i)
}

func TestCordon(t *testing.T) {
	t.Parallel()

	nodes := newNodeListFixture(nlConfig{count: 1, cpu: []string{"2.0"}})
	cs := &fake.Clientset{}
	s, err := NewServer("", "app=game-server", "0.5", 5, time.Second)
	assert.Nil(t, err)
	s.cs = cs
	defaultListNodeReactor(cs, nodes)
	defaultUpdateNodeReactor(cs, nodes)

	now := time.Now().UTC()
	node := nodes.Items[0]
	err = s.cordon(&node, true)
	assert.Nil(t, err)
	assert.True(t, node.Spec.Unschedulable)
	var ts time.Time
	err = ts.UnmarshalText([]byte(node.ObjectMeta.Annotations[timestampAnnotation]))
	assert.Nil(t, err)
	assert.True(t, ts.Equal(now) || ts.After(now), "Now: %v is not equal to or after %v", now, ts)

	nl, err := s.newNodeList()
	assert.Nil(t, err)
	assert.True(t, nl.nodes.Items[0].Spec.Unschedulable)
	assert.Equal(t, 0, len(nl.availableNodes()))

	err = s.cordon(&node, false)
	assert.Nil(t, err)
	assert.False(t, node.Spec.Unschedulable)
	err = ts.UnmarshalText([]byte(node.ObjectMeta.Annotations[timestampAnnotation]))
	assert.Nil(t, err)
	assert.True(t, ts.Equal(now) || ts.After(now), "Now: %v is not equal to or after %v", now, ts)

	nl, err = s.newNodeList()
	assert.Nil(t, err)
	assert.False(t, nl.nodes.Items[0].Spec.Unschedulable)
	assert.Equal(t, 1, len(nl.availableNodes()))
}
