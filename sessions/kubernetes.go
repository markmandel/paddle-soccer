package sessions

import (
	"fmt"

	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

// ClientSet get the kubernetes clientset
func ClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return cs, nil
}

// HostNameAndIP returns a map of hostname (key) to external IP (value)
func HostNameAndIP(cs *kubernetes.Clientset) (map[string]string, error) {
	result := make(map[string]string)
	var list *v1.NodeList
	list, err := cs.CoreV1().Nodes().List(v1.ListOptions{})

	if err != nil {
		return result, err
	}

	for _, i := range list.Items {
		k := i.ObjectMeta.Name
		var addr string
		for _, a := range i.Status.Addresses {
			if a.Type == v1.NodeExternalIP {
				addr = a.Address
				break
			}
		}

		if addr == "" {
			return result, fmt.Errorf("[Error][Kubernetes] Could not find an external ip for Node: #%v", i)
		}

		result[k] = addr
	}

	log.Printf("[Info][Kubernetes] List of nodes and ips: %#v", result)

	return result, nil
}

// ExternalNodeIPofPod gets the external IP address of the node a pod is on
func ExternalNodeIPofPod(cs *kubernetes.Clientset, s Session, m map[string]string) (string, error) {
	log.Printf("[Info][Kubernetes] Retrieving pod information for pod: %v", s.ID)

	var pod *v1.Pod
	pod, err := cs.CoreV1().Pods("default").Get(s.ID)

	if err != nil {
		log.Printf("[Error][Register] Error getting pod information for pod %v: %v", s.ID, err)
		return "", err
	}

	nodeName := pod.Spec.NodeName
	return m[nodeName], nil
}
