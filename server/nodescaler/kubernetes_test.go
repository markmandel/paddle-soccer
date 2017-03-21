package nodescaler

import (
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
	core "k8s.io/client-go/testing"
)

func TestListNodePods(t *testing.T) {
	fixture := &v1.PodList{Items: []v1.Pod{
		{ObjectMeta: v1.ObjectMeta{Name: "pod1", Namespace: "default"}},
		{ObjectMeta: v1.ObjectMeta{Name: "pod2", Namespace: "default"}}}}

	s, err := NewServer("", "app=game-server", "0.5", 5)
	assert.Nil(t, err)
	sc := fake.NewSimpleClientset(fixture)
	s.cs = sc

	pods, err := s.listNodePods(v1.Node{})
	assert.Nil(t, err)
	assert.Equal(t, fixture, pods)
}

func TestSumCPUResourceRequests(t *testing.T) {
	fixture := &v1.PodList{Items: []v1.Pod{
		{ObjectMeta: v1.ObjectMeta{Name: "pod1", Namespace: "default"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("0.5")}}}}}},
		{ObjectMeta: v1.ObjectMeta{Name: "pod2", Namespace: "default"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("0.3")}}}}}}}}

	expected := resource.MustParse("0.8")
	assert.Equal(t, expected.MilliValue(), sumCPUResourceRequests(fixture))
}

func TestNewNodeList(t *testing.T) {
	nodes := &v1.NodeList{Items: []v1.Node{
		{ObjectMeta: v1.ObjectMeta{Name: "foo", Labels: map[string]string{"app": "game-server"}},
			Status: v1.NodeStatus{Capacity: v1.ResourceList{v1.ResourceCPU: resource.MustParse("4.0")}}},
		{ObjectMeta: v1.ObjectMeta{Name: "bar", Labels: map[string]string{"app": "game-server"}},
			Status: v1.NodeStatus{Capacity: v1.ResourceList{v1.ResourceCPU: resource.MustParse("4.0")}}}}}

	pl1 := &v1.PodList{Items: []v1.Pod{
		{ObjectMeta: v1.ObjectMeta{Name: "pod1", Namespace: "default"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("0.5")}}}}}},
		{ObjectMeta: v1.ObjectMeta{Name: "pod2", Namespace: "default"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("0.3")}}}}}}}}

	pl2 := &v1.PodList{Items: []v1.Pod{
		{ObjectMeta: v1.ObjectMeta{Name: "pod3", Namespace: "default"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("1.2")}}}}}}}}

	cs := &fake.Clientset{}

	cs.AddReactor("list", "nodes", func(a core.Action) (bool, runtime.Object, error) {
		log.Print("Asking for list of nodes...")
		return true, nodes, nil
	})

	cs.AddReactor("list", "pods", func(a core.Action) (bool, runtime.Object, error) {
		var obj *v1.PodList
		la := a.(core.ListAction)

		switch la.GetListRestrictions().Fields.String() {
		case "spec.nodeName=foo":
			obj = pl1
		case "spec.nodeName=bar":
			obj = pl2
		}
		return true, obj, nil
	})

	s, err := NewServer("", "app=game-server", "0.5", 5)
	assert.Nil(t, err)
	s.cs = cs

	nodeList, err := s.newNodeList()
	assert.Nil(t, err)
	assert.Equal(t, nodes, nodeList.nodes)
	assert.Equal(t, pl1, nodeList.nodePods(nodes.Items[0]))
	assert.Equal(t, pl2, nodeList.nodePods(nodes.Items[1]))
}

func TestCpuRequestsAvailable(t *testing.T) {
	nodes := &v1.NodeList{Items: []v1.Node{
		{ObjectMeta: v1.ObjectMeta{Name: "foo", Labels: map[string]string{"app": "game-server"}},
			Status: v1.NodeStatus{Capacity: v1.ResourceList{v1.ResourceCPU: resource.MustParse("2.0")}}},
		{ObjectMeta: v1.ObjectMeta{Name: "bar", Labels: map[string]string{"app": "game-server"}},
			Status: v1.NodeStatus{Capacity: v1.ResourceList{v1.ResourceCPU: resource.MustParse("2.0")}}}}}

	pl1 := &v1.PodList{Items: []v1.Pod{
		{ObjectMeta: v1.ObjectMeta{Name: "pod1", Namespace: "default"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("0.5")}}}}}},
		{ObjectMeta: v1.ObjectMeta{Name: "pod2", Namespace: "default"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("0.3")}}}}}}}}

	pl2 := &v1.PodList{Items: []v1.Pod{
		{ObjectMeta: v1.ObjectMeta{Name: "pod3", Namespace: "default"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("1.8")}}}}}}}}

	cs := &fake.Clientset{}

	cs.AddReactor("list", "nodes", func(a core.Action) (bool, runtime.Object, error) {
		log.Print("Asking for list of nodes...")
		return true, nodes, nil
	})

	cs.AddReactor("list", "pods", func(a core.Action) (bool, runtime.Object, error) {
		var obj *v1.PodList
		la := a.(core.ListAction)

		switch la.GetListRestrictions().Fields.String() {
		case "spec.nodeName=foo":
			obj = pl1
		case "spec.nodeName=bar":
			obj = pl2
		}
		return true, obj, nil
	})

	s, err := NewServer("", "app=game-server", "0.5", 5)
	assert.Nil(t, err)
	s.cs = cs

	nl, err := s.newNodeList()
	assert.Nil(t, err)

	count := s.cpuRequestsAvailable(nl)
	assert.Equal(t, int64(2), count)
}
