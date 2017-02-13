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

package sessions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
)

func TestServerHostNameAndIP(t *testing.T) {
	nl := v1.NodeList{Items: []v1.Node{{
		ObjectMeta: v1.ObjectMeta{Name: "foo"},
		Status: v1.NodeStatus{
			Addresses: []v1.NodeAddress{{Type: v1.NodeExternalIP, Address: "1.2.3.4"}},
		},
	}}}

	s := NewServer("", "", "")
	sc := fake.NewSimpleClientset(&nl)
	s.cs = sc

	result, err := s.hostNameAndIP()
	assert.Nil(t, err)
	expected := map[string]string{"foo": "1.2.3.4"}
	assert.Equal(t, expected, result)
}

func TestServerExternalNodeIPofPod(t *testing.T) {
	p := v1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			NodeName: "foo",
		},
	}

	s := NewServer("", "", "")
	sc := fake.NewSimpleClientset(&p)
	s.cs = sc

	fixture := map[string]string{"foo": "1.2.3.4"}
	ip, err := s.externalNodeIPofPod(Session{}, fixture)

	assert.Nil(t, err)
	assert.Equal(t, fixture["foo"], ip)
}

func TestServerCreateSessionPod(t *testing.T) {
	s := NewServer("", "", "foo")
	sc := fake.NewSimpleClientset()
	s.cs = sc

	_, err := s.createSessionPod()
	assert.Nil(t, err)

	list, err := sc.CoreV1().Pods("default").List(v1.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, list.Items, 1)
}
