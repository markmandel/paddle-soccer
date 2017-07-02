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
	"math"
	"sort"

	"k8s.io/client-go/pkg/api/v1"
)

// NodePool is an interface to
// being able to manipulate a pool for nodes
// for Kubernetes.
type NodePool interface {
	// Increase the node pool to a given size.
	// Should ignore requests to make the nodepool smaller
	IncreaseToSize(int64) error
}

// scale scales nodes up and down, depending on CPU constraints
func (s Server) scaleNodes() error {
	nl, err := s.newNodeList()
	if err != nil {
		return err
	}

	available := nl.cpuRequestsAvailable()
	log.Printf("[Info][scaleNodes] CPU Requests blocks of %vm. Available: %v. Requires a buffer of %v", s.cpuRequest, available, s.bufferCount)
	if available < s.bufferCount {
		finished, err := s.uncordonNodes(nl, s.bufferCount-available)
		// short circuit if uncordoning means we have enough buffer now
		if err != nil || finished {
			return err
		}

		nl, err := s.newNodeList()
		if err != nil {
			return err
		}
		// recalculate
		available = nl.cpuRequestsAvailable()
		return s.increaseNodes(nl, s.bufferCount-available)
	} else if s.bufferCount < available {
		return s.cordonNodes(nl, available-s.bufferCount)
	}

	return nil
}

// increaseNodes increases the size of the managed nodepool
// by a given number of cpu blocks
func (s Server) increaseNodes(nl *nodeList, gameNumber int64) error {
	if gameNumber <= 0 {
		log.Printf("[Warn][IncreaseNodes] Cannot increase nodes by a 0 or negative number. %v", gameNumber)
		return nil
	}

	log.Printf("[Info][increaseNodes] Attempting to increase nodelist of %v, by %v cpu blocks", len(nl.nodes.Items), gameNumber)

	// determine how many nodes we need to be equal to or greater
	// than the bufferCount requirements.

	// this won't be totally accurate, since there are a few more containers on there
	// but it will do for now, and if not, the next round will pull it up
	capacity := nl.nodes.Items[0].Status.Capacity[v1.ResourceCPU]
	cpuRequest := gameNumber * s.cpuRequest
	diff := int64(math.Ceil(float64(cpuRequest) / float64(capacity.MilliValue())))
	log.Printf("[Info][increaseNodes] Adding %v nodes to the node pool", diff)

	return s.nodePool.IncreaseToSize(int64(len(nl.availableNodes())) + diff)
}

// uncordonNodes searches through all the available nodes and uncordons
// those that are most at capacity up to the gameNumber limit
// returns true if it successfully fullfilled the gameNumber request that was passed in
func (s Server) uncordonNodes(nl *nodeList, gameNumber int64) (bool, error) {
	if gameNumber <= 0 {
		log.Printf("[Warn][uncordonNodes] Cannot uncordon nodes by a 0 or negative number. %v", gameNumber)
		return true, nil
	}

	// first we need to take all nodes, and sort by decreasing allocatable value
	nodes := nl.cordonedNodes()
	if len(nodes) == 0 {
		log.Print("[Info][uncordonNodes] No nodes that are unscheduled, exiting...")
		return false, nil
	}

	// sort by the node that has the most pods on it currently (least availability)
	sort.Slice(nodes, func(i, j int) bool {
		iSum := nl.sumResourceLimit(nodes[i], v1.ResourceCPU)
		jSum := nl.sumResourceLimit(nodes[j], v1.ResourceCPU)

		return jSum < iSum
	})

	cpuRequest := s.cpuRequest * gameNumber

	log.Printf("[Info][uncordonNodes] Uncordoning nodes. Requesting %v game servers for a cpuRequest of %v", gameNumber, cpuRequest)

	for _, n := range nodes {
		capacity := n.Status.Capacity[v1.ResourceCPU]
		err := s.cordon(&n, false)
		if err != nil {
			return false, err
		}

		available := capacity.MilliValue() - nl.sumResourceLimit(n, v1.ResourceCPU)
		cpuRequest = cpuRequest - available
		log.Printf("[Info][uncordonNodes] %v cpuRequest remaining after uncordoning node, and adding %v cpu", cpuRequest, available)

		if cpuRequest <= 0 {
			return true, nil
		}
	}
	return false, nil
}

// cordonNodes decrease the number of available nodes by the given number of cpu blocks (but not over),
// but cordoning those nodes that have the least number of games currently on them
func (s Server) cordonNodes(nl *nodeList, gameNumber int64) error {
	if gameNumber <= 0 {
		log.Printf("[Warn][CordonNodes] Cannot cordon nodes by a 0 or negative number. %v", gameNumber)
		return nil
	}

	// how many nodes (n) do we have to delete such that we are cordoning no more
	// than the gameNumber
	capacity := nl.nodes.Items[0].Status.Capacity[v1.ResourceCPU] //assuming all nodes are the same
	cpuRequest := gameNumber * s.cpuRequest
	diff := int64(math.Floor(float64(cpuRequest) / float64(capacity.MilliValue())))

	if diff <= 0 {
		log.Print("[Info][CordonNodes] No nodes to be cordoned.")
		return nil
	}

	log.Printf("[Info][CordonNodes] Cordoning %v nodes", diff)

	// sort the nodes, such that the one with the least number of games are first
	nodes := nl.nodes.Items
	sort.Slice(nodes, func(i, j int) bool {
		return len(nl.nodePods(nodes[i]).Items) < len(nl.nodePods(nodes[j]).Items)
	})

	// grab the first n number of them
	cNodes := nodes[0:diff]

	// cordon them all
	for _, n := range cNodes {
		log.Printf("[Info][CordonNodes] Cordoning node: %v", n.Name)
		err := s.cordon(&n, true)
		log.Printf("[Debug][CordonNodes] Status: %v, %v", n.Name, n.Spec.Unschedulable)
		if err != nil {
			return err
		}
	}

	return nil
}
