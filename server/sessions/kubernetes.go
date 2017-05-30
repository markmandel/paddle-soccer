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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// hostNameAndIP returns a map of Kubernetes node hostname (key) to external IP (value)
func (s *Server) hostNameAndIP() (map[string]string, error) {
	result := make(map[string]string)
	list, err := s.cs.CoreV1().Nodes().List(metav1.ListOptions{})

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

	pod, err := s.cs.CoreV1().Pods("default").Get(sess.ID, metav1.GetOptions{})
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

	labels := map[string]string{
		"sessions": "game",
	}

	pod := v1.Pod{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "sessions-game-",
			Labels:       labels,
			Namespace:    namespace,
		},
		Spec: v1.PodSpec{
			HostNetwork:   true,
			RestartPolicy: v1.RestartPolicyNever,
			NodeSelector:  s.gameNodeSelector,
			Containers: []v1.Container{
				{
					Name:            "sessions-game",
					Image:           s.gameServerImage,
					ImagePullPolicy: v1.PullAlways, // this is just for dev / POC
					Resources:       v1.ResourceRequirements{Limits: v1.ResourceList{v1.ResourceCPU: s.cpuLimit}},
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
			// Make it so that each game server is as close to each other as they can be
			// This lowers fragmentation, and makes it easier to scale down
			Affinity: &v1.Affinity{
				PodAffinity: &v1.PodAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
						{
							Weight: 1, PodAffinityTerm: v1.PodAffinityTerm{
								LabelSelector: &metav1.LabelSelector{MatchLabels: labels},
								TopologyKey:   metav1.LabelHostname,
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
