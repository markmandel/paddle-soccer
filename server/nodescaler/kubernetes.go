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
	"sync"
	"time"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	// Attribute key for where to put the timestamp when
	// making changes to a node
	timestampAnnotation = "nodescale/timestamp"
)

// nodeList is the set of current nodes that this
// server manages, with the accompanying pods for
// each node
type nodeList struct {
	nodes      *v1.NodeList
	pods       map[string]*v1.PodList
	cpuRequest int64
}

// newNodeList queries kubernetes to get a list of nodes
// with their associated pods
func (s *Server) newNodeList() (*nodeList, error) {
	var result *nodeList

	nodes, err := s.cs.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: s.nodeSelector})
	if err != nil {
		return result, errors.Wrap(err, "Could not get node list from Kubernetes")
	}

	result = &nodeList{nodes: nodes, pods: map[string]*v1.PodList{}, cpuRequest: s.cpuRequest}
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

// availableNodes returns nodes that are available to be used.
// This could mean they are ready
// that they are scheduled.
func (nl nodeList) availableNodes() []v1.Node {
	var result []v1.Node
	for _, n := range nl.nodes.Items {
		if nodeReady(n) && !n.Spec.Unschedulable {
			result = append(result, n)
		}
	}
	return result
}

// cordonedNodes returns all notes that are set to being
// unscheduled
func (nl nodeList) cordonedNodes() []v1.Node {
	var result []v1.Node
	for _, n := range nl.nodes.Items {
		if n.Spec.Unschedulable {
			result = append(result, n)
		}
	}
	return result
}

// cpuRequestsAvailable looks at each node's cpu capacity,
// the current available free space in it, and determines
// how many cpu requests are can be be fit in the
// remaining space on each node
func (nl *nodeList) cpuRequestsAvailable() int64 {
	var avail int64
	for _, n := range nl.availableNodes() {
		capacity := n.Status.Capacity[v1.ResourceCPU]
		requests := nl.sumResourceLimit(n, v1.ResourceCPU)
		diff := capacity.MilliValue() - requests

		// take advantage of the fact we are using
		// int64's and remainders / fractions are
		// totally not what we want
		avail += diff / nl.cpuRequest
	}

	return avail
}

// sumResourceLimit returns the sum of all of the limits for given resource for all the pods in a
// specific Node
func (nl *nodeList) sumResourceLimit(n v1.Node, name v1.ResourceName) int64 {
	var total int64
	pl := nl.nodePods(n)
	for _, p := range pl.Items {
		for _, c := range p.Spec.Containers {
			r := c.Resources.Limits[name]
			total += r.MilliValue()
		}
	}

	return total
}

// gameWatcher implements Kubernetes watch.Interface to allow for
// scaling up to be processed whenever a game event occurs
// Adds a value to the events channel whenever a game is
// added to the nodepool in question, or deleted
type gameWatcher struct {
	watcher watch.Interface
	events  chan bool
	// Wait Group to ensure that closing of channels on stop
	// doesn't interrupt a currently processing event.
	wg sync.WaitGroup
}

// newGameWatcher returns a new game watcher
func (s *Server) newGameWatcher() (*gameWatcher, error) {
	g := &gameWatcher{events: make(chan bool)}

	watcher, err := s.cs.CoreV1().Pods(api.NamespaceAll).Watch(metav1.ListOptions{LabelSelector: "sessions=game"})
	if err != nil {
		return g, errors.Wrapf(err, "Error watching pods for selector: %#v", s.nodeSelector)
	}
	g.watcher = watcher

	return g, nil
}

// start starts the game watcher, watching the K8 event stream
func (g *gameWatcher) start() {
	go func() {
		// WaitGroup for ensuring that if we are shutting down
		// we don't shut down the events/deleted channels
		// before ResultChan's events are fully processed.
		g.wg.Add(1)
		for e := range g.watcher.ResultChan() {
			if e.Type == watch.Added || e.Type == watch.Deleted {
				g.events <- true
			}
		}
		g.wg.Done()
	}()
}

// stop closes all the channels, when you are done
func (g *gameWatcher) stop() {
	g.watcher.Stop()
	g.wg.Wait()
	close(g.events)
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

	pods, err := s.cs.CoreV1().Pods(api.NamespaceAll).List(metav1.ListOptions{FieldSelector: fs.String()})
	return pods, errors.Wrapf(err, "Could not get pods for Node: %v", n.Name)
}

// cordon sets or unsets a node to being unschedulable
// a 'true' parameter will set a node to being unschedulable (cordoned)
// this also sets a timestamp annotation on the node, to track when this was
// last done.
func (s *Server) cordon(n *v1.Node, unscheduled bool) error {
	now, err := time.Now().UTC().MarshalText()
	if err != nil {
		return errors.Wrap(err, "Could not marshall now datetime as string")
	}

	n.Spec.Unschedulable = unscheduled
	n.ObjectMeta.Annotations[timestampAnnotation] = string(now)
	_, err = s.cs.CoreV1().Nodes().Update(n)
	return errors.Wrapf(err, "Error Updating Node %#v", n)
}
