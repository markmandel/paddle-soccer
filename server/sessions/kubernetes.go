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
	"log"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

// clientSet returns the production Kubernetes clientSet
func clientSet() (kubernetes.Interface, error) {
	log.Print("[Info][Kubernetes] Connecting to Kubernetes API...")
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to Kubernetes API")
	}
	log.Print("[Info][Kubernetes] Connected to Kubernetes API")
	return kubernetes.NewForConfig(config)
}

// hostNameAndIP returns a map of Kubernetes node hostname (key) to external IP (value)
func (s *Server) hostNameAndIP() (map[string]string, error) {
	result := make(map[string]string)
	list, err := s.cs.CoreV1().Nodes().List(v1.ListOptions{})

	if err != nil {
		return result, errors.Wrap(err, "Error Listing nodes")
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
			return result, errors.Errorf("Could not find an external ip for Node: #%v", i)
		}

		result[k] = addr
	}

	log.Printf("[Info][Kubernetes] List of nodes and ips: %#v", result)

	return result, nil
}

// externalNodeIPofPod gets the external IP address of the node a pod is on
func (s *Server) externalNodeIPofPod(sess Session, m map[string]string) (string, error) {
	log.Printf("[Info][Kubernetes] Retrieving pod information for pod: %v", sess.ID)

	pod, err := s.cs.CoreV1().Pods("default").Get(sess.ID)
	if err != nil {
		return "", errors.Wrapf(err, "Error getting pod information for pod %v", sess.ID)
	}
	nodeName := pod.Spec.NodeName
	return m[nodeName], nil
}

// createSessionPod creates a pod containing the game server, for a game session
func (s *Server) createSessionPod() (string, error) {
	log.Print("[Info][create] Creating new game session")
	namespace := "default"

	pod := v1.Pod{
		TypeMeta: unversioned.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: v1.ObjectMeta{
			GenerateName: "sessions-game-",
			Labels: map[string]string{
				"sessions": "game",
			},
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			HostNetwork:   true,
			RestartPolicy: v1.RestartPolicyNever,
			NodeSelector:  s.gameNodeSelector,
			Containers: []v1.Container{
				{
					Name:            "sessions-game",
					Image:           s.gameServerImage,
					ImagePullPolicy: v1.PullAlways, // TODO: make this an env var / this is just for dev
					Env: []v1.EnvVar{
						{
							Name: "SESSION_NAME",
							ValueFrom: &v1.EnvVarSource{
								FieldRef: &v1.ObjectFieldSelector{
									FieldPath: "metadata.name",
								},
							},
						},
					},
				},
			},
		},
	}

	log.Printf("[Info][Kubernetes] Creating pod: %#v", pod)
	result, err := s.cs.CoreV1().Pods(namespace).Create(&pod)
	if err != nil {
		return "", errors.Wrap(err, "Error creating pod")
	}
	log.Printf("[Info][Kubernetes] Created pod: %v", result.Name)
	return result.Name, nil
}
