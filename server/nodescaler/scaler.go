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

	available := s.cpuRequestsAvailable(nl)
	log.Printf("[Info][scaleNodes] CPU Requests blocks of %vm. Available: %v, with buffer of %v", s.cpuRequest, available, s.bufferCount)
	if available < s.bufferCount {
		return s.increaseNodes(nl, s.bufferCount-available)
	}

	return nil
}

// increaseNodes increases the size of the managed nodepool
// by a given amount of cpu
func (s Server) increaseNodes(nl *nodeList, cpuAmount int64) error {
	if cpuAmount <= 0 {
		log.Printf("[Warn][IncreaseNodes] Cannot increase nodes by a 0 or negative number. %v", cpuAmount)
		return nil
	}

	log.Printf("[Info][increaseNodes] Attempting to increase nodelist of %v, by %v cpu blocks", len(nl.nodes.Items), cpuAmount)

	// determine how many nodes we need to be equal to or greater
	// than the bufferCount requirements.

	// this won't be totally accurate, since there are a few more containers on there
	// but it will do for now, and if not, the next round will pull it up
	nodeCapacity := nl.nodes.Items[0].Status.Capacity[v1.ResourceCPU]
	cpuRequest := cpuAmount * s.cpuRequest
	diff := int64(math.Ceil(float64(cpuRequest) / float64(nodeCapacity.MilliValue())))
	log.Printf("[Info][increaseNodes] Adding %v nodes to the node pool", diff)

	return s.nodePool.IncreaseToSize(int64(len(nl.availableNodes())) + diff)
}
