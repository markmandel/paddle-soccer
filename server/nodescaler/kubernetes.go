package nodescaler

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
)

// nodeList is the set of current nodes that this
// server manages, with the accompanying pods for
// each node
type nodeList struct {
	nodes *v1.NodeList
	pods  map[string]*v1.PodList
}

// nodePods retrieves a pod for a given node
func (nl nodeList) nodePods(n v1.Node) *v1.PodList {
	return nl.pods[n.Name]
}

// cpuRequestsAvailable looks at each node's cpu capacity,
// the current available free space in it, and determines
// how many cpu requests are can be be fit in the
// remaining space on each node
func (s *Server) cpuRequestsAvailable(nl *nodeList) int64 {
	var avail int64
	for _, n := range nl.nodes.Items {
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
