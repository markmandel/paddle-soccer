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
	"strconv"
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	core "k8s.io/client-go/testing"
)

// readyNodeCondition for when a node is ready, as most in this test suite will be
var readyNodeCondition = []v1.NodeCondition{{Type: v1.NodeReady, Status: v1.ConditionTrue}}

// defaultUpdateNodeReactor adds the standard reactor function for updating
// a nodelist on pod update
func defaultUpdateNodeReactor(cs *fake.Clientset, target *v1.NodeList) {
	cs.AddReactor("update", "nodes", func(a core.Action) (bool, runtime.Object, error) {
		ua := a.(core.UpdateAction)
		n := ua.GetObject().(*v1.Node)

		for i, ns := range target.Items {
			if ns.Name == n.Name {
				log.Printf("[Debug][UpdateNode] Updating node: %#v", n)
				target.Items[i] = *n
			}
		}

		return true, n, nil
	})
}

// defaultListNodeReactor add the default list-nodes reactors to return the
// given target
func defaultListNodeReactor(cs *fake.Clientset, target *v1.NodeList) {
	cs.AddReactor("list", "nodes", func(a core.Action) (bool, runtime.Object, error) {
		return true, target, nil
	})
}

type nlConfig struct {
	count int
	cpu   []string
}

// nodeListBuilder builds out standard nodeLists, with the cpu resource values
// set
func newNodeListFixture(config nlConfig) *v1.NodeList {
	result := &v1.NodeList{}
	for i := 0; i < config.count; i++ {
		n := v1.Node{ObjectMeta: metav1.ObjectMeta{Name: "node" + strconv.Itoa(i),
			Labels: map[string]string{"app": "game-server"}, Annotations: map[string]string{}},
			Spec: v1.NodeSpec{Unschedulable: false},
			Status: v1.NodeStatus{Capacity: v1.ResourceList{},
				Conditions: readyNodeCondition}}

		if len(config.cpu) > 0 {
			n.Status.Capacity[v1.ResourceCPU] = resource.MustParse(config.cpu[i])
		}

		result.Items = append(result.Items, n)
	}

	return result
}

// assertAllUnscheduled checks all nodes are unscheduled
func assertAllUnscheduled(t *testing.T, nodes *v1.NodeList) {
	for _, n := range nodes.Items {
		assert.False(t, n.Spec.Unschedulable, "Node %v, should not be schedulable", n.Name)
	}
}

// newPodListFixture returns a standard pod list, for each cpu resource
// value provided.
func newPodListFixture(cpu []string) *v1.PodList {
	result := &v1.PodList{}

	for i, c := range cpu {
		pod := v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod" + strconv.Itoa(i), Namespace: "default"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{v1.ResourceCPU: resource.MustParse(c)}}}}}}

		result.Items = append(result.Items, pod)
	}

	return result
}

// defaultListPodReactor create a standard pod list reactor, where the node name
// matches the podlist associated
func defaultListPodReactor(cs *fake.Clientset, np map[string]*v1.PodList) {
	cs.AddReactor("list", "pods", func(a core.Action) (bool, runtime.Object, error) {
		la := a.(core.ListAction)
		p := strings.Split(la.GetListRestrictions().Fields.String(), "=")
		nn := p[1]

		pods, ok := np[nn]
		if ok {
			return true, pods, nil
		}

		return false, nil, nil
	})
}
