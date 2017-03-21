package kube

import (
	"log"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ClientSet returns the production Kubernetes clientSet
func ClientSet() (kubernetes.Interface, error) {
	log.Print("[Info][Kubernetes] Connecting to Kubernetes API...")
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to Kubernetes API")
	}
	log.Print("[Info][Kubernetes] Connected to Kubernetes API")
	return kubernetes.NewForConfig(config)
}
