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
	"github.com/pkg/errors"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/pkg/watch"
)

// nodeList is the set of current nodes that this
// server manages, with the accompanying pods for
// each node
type nodeList struct {
	nodes *v1.NodeList
	pods  map[string]*v1.PodList
}

// newNodeList queries kubernetes to get a list of nodes
// with their associated pods
func (s *Server) newNodeList() (*nodeList, error) {
	var result *nodeList

	nodes, err := s.cs.CoreV1().Nodes().List(v1.ListOptions{LabelSelector: s.nodeSelector})
	if err != nil {
		return result, errors.Wrap(err, "Could not get node list from Kubernetes")
	}

	result = &nodeList{nodes: nodes, pods: map[string]*v1.PodList{}}
	for _, n := range nodes.Items {
		pods, err := s.listNodePods(n)
		if err != nil {
			return result, err
		}
		result.pods[n.Name] = pods
	}
	return result, err
}

// nodePods retrieves a pod for a given node
func (nl nodeList) nodePods(n v1.Node) *v1.PodList {
	return nl.pods[n.Name]
}

// returns nodes that are available to be used.
// This could mean they are ready, or (TBD)
// that they are available to be scheduled.
func (nl nodeList) availableNodes() []v1.Node {
	result := []v1.Node{}
	for _, n := range nl.nodes.Items {
		if nodeReady(n) {
			result = append(result, n)
		}
	}

	return result
}

// implements Kubernetes watch.Interface to allow for
// scaling up to be processed whenever a game event occurs
// Adds a value to the event channel whenever a game is
// added to the nodepool in question
type gameWatcher struct {
	watcher watch.Interface
	event   chan bool
}

// newGameWatcher returns a new game watcher
func (s *Server) newGameWatcher() (*gameWatcher, error) {
	g := &gameWatcher{event: make(chan bool)}

	watcher, err := s.cs.CoreV1().Pods(api.NamespaceAll).Watch(v1.ListOptions{LabelSelector: "sessions=game"})
	if err != nil {
		return g, errors.Wrapf(err, "Error watching pods for selector: %#v", s.nodeSelector)
	}
	g.watcher = watcher

	return g, nil
}

// start starts the game watcher, watching the K8 event stream
func (g *gameWatcher) start() {
	go func() {
		for e := range g.watcher.ResultChan() {
			if e.Type == watch.Added {
				g.event <- true
			}
		}
	}()
}

// stop closes all the channels, when you are done
func (g *gameWatcher) stop() {
	g.watcher.Stop()
	close(g.event)
}

// cpuRequestsAvailable looks at each node's cpu capacity,
// the current available free space in it, and determines
// how many cpu requests are can be be fit in the
// remaining space on each node
func (s *Server) cpuRequestsAvailable(nl *nodeList) int64 {
	var avail int64
	for _, n := range nl.availableNodes() {
		capacity := n.Status.Capacity[v1.ResourceCPU]
		requests := sumCPUResourceRequests(nl.nodePods(n))
		diff := capacity.MilliValue() - requests

		// take advantage of the fact we are using
		// int64's and remainders / fractions are
		// totally not what we want
		avail += diff / s.cpuRequest
	}

	return avail
}

// nodeReady check if a node's kublet is ready to work
func nodeReady(n v1.Node) bool {
	for _, c := range n.Status.Conditions {
		if c.Type == v1.NodeReady && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

// listNodsPods list all the pods for this node
func (s *Server) listNodePods(n v1.Node) (*v1.PodList, error) {
	fs, err := fields.ParseSelector("spec.nodeName=" + n.Name)
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse selector")
	}

	pods, err := s.cs.CoreV1().Pods(api.NamespaceAll).List(v1.ListOptions{FieldSelector: fs.String()})
	return pods, errors.Wrapf(err, "Could not get pods for Node: %v", n.Name)
}

// sumCPUResourceRequests returns the sum of all the pod
// CPU resource requests
func sumCPUResourceRequests(pl *v1.PodList) int64 {
	var total int64
	for _, p := range pl.Items {
		for _, c := range p.Spec.Containers {
			r := c.Resources.Requests[v1.ResourceCPU]
			total += r.MilliValue()
		}
	}

	return total
}
