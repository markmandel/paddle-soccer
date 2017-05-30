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
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	core "k8s.io/client-go/testing"
)

func TestScaleNodes(t *testing.T) {
	nodes := &v1.NodeList{Items: []v1.Node{
		{ObjectMeta: metav1.ObjectMeta{Name: "foo", Labels: map[string]string{"app": "game-server"}},
			Status: v1.NodeStatus{Capacity: v1.ResourceList{v1.ResourceCPU: resource.MustParse("2.0")},
				Conditions: readyNodeCondition}},
		{ObjectMeta: metav1.ObjectMeta{Name: "bar", Labels: map[string]string{"app": "game-server"}},
			Status: v1.NodeStatus{Capacity: v1.ResourceList{v1.ResourceCPU: resource.MustParse("2.0")},
				Conditions: readyNodeCondition}}}}

	cs := &fake.Clientset{}
	cs.AddReactor("list", "nodes", func(a core.Action) (bool, runtime.Object, error) {
		return true, nodes, nil
	})

	s, err := NewServer("", "app=game-server", "0.5", 5)
	assert.Nil(t, err)
	s.cs = cs

	expected := int64(len(nodes.Items))
	mock := &NodePoolMock{size: expected}
	s.nodePool = mock
	err = s.scaleNodes()
	assert.Nil(t, err)
	assert.Equal(t, expected, mock.size)

	cs.AddReactor("list", "pods", func(a core.Action) (bool, runtime.Object, error) {
		if a.(core.ListAction).GetListRestrictions().Fields.String() == "spec.nodeName=foo" {
			return true,
				&v1.PodList{Items: []v1.Pod{
					{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "default"},
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("0.5")}}}}}},
					{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "default"},
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("0.3")}}}}}}}}, nil
		}
		return false, nil, nil
	})

	s.nodePool = mock
	err = s.scaleNodes()
	assert.Nil(t, err)
	assert.Equal(t, expected, mock.size)

	cs.AddReactor("list", "pods", func(a core.Action) (bool, runtime.Object, error) {
		if a.(core.ListAction).GetListRestrictions().Fields.String() == "spec.nodeName=bar" {
			return true,
				&v1.PodList{Items: []v1.Pod{
					{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "default"},
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("1.8")}}}}}}}}, nil
		}
		return false, nil, nil
	})

	s.nodePool = mock
	err = s.scaleNodes()
	assert.Nil(t, err)
	assert.Equal(t, int64(3), mock.size)
}

type NodePoolMock struct {
	size int64
}

// IncreaseToSize
func (npm *NodePoolMock) IncreaseToSize(size int64) error {
	if size <= npm.size {
		log.Printf("[Debug][Mock:IncreaseToSize] Ignoring resize to %v, as size is already %v", size, npm.size)
		return nil
	}
	log.Printf("[Debug][Mock:IncreaseToSize] Resising to: %v", size)
	npm.size = size
	return nil
}
