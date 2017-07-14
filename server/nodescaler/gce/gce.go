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

package gce

import (
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	"k8s.io/client-go/pkg/api/v1"
)

type NodePool struct {
	project              string
	zone                 string
	instanceGroup        string
	instanceGroupService *compute.InstanceGroupManagersService
}

// NewNodePool returns a new gce nodepool
func NewNodePool(n v1.Node) (*NodePool, error) {
	np := &NodePool{}
	var err error
	np.project, err = metadata.ProjectID()
	if err != nil {
		return np, errors.Wrap(err, "Could not get project id from metadata")
	}

	np.zone, err = metadata.Zone()
	if err != nil {
		return np, errors.Wrap(err, "Could not get zone from metadata")
	}

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope, container.CloudPlatformScope)
	if err != nil {
		return np, errors.Wrap(err, "Could not create google default client")
	}

	ig, err := np.nodeInstanceGroup(n, client)
	if err != nil {
		return np, err
	}
	np.instanceGroup = ig

	cmpService, err := compute.New(client)
	if err != nil {
		return np, errors.Wrap(err, "Could not create compute service")
	}
	np.instanceGroupService = compute.NewInstanceGroupManagersService(cmpService)

	log.Printf("[Info][NewNodePool] Connected to nodepool: %#v", np)

	return np, nil
}

// Increases the nodepool size, if the size is greater than the current
// target size for the underlying instance group
func (np *NodePool) IncreaseToSize(size int64) error {
	log.Printf("[Info][IncreaseToSize] Attempting to increase instance group %v to size %v", np.instanceGroup, size)
	ig, err := np.instanceGroupService.Get(np.project, np.zone, np.instanceGroup).Do()
	if err != nil {
		return errors.Wrapf(err, "Could not get instance group data for %v", np.instanceGroup)
	}

	if size <= ig.TargetSize {
		log.Printf("[Info][IncreaseToSize] Ignoring change to %v's size, as target size is already %v",
			np.instanceGroup, ig.TargetSize)
	}

	_, err = np.instanceGroupService.Resize(np.project, np.zone, np.instanceGroup, size).Do()
	if err != nil {
		return errors.Wrapf(err, "Error resizing %v to size %v", np.instanceGroup, size)
	}

	return nil
}

// nodeInstanceGroup works out the underlying instance group
// that this K8s node is managed by
func (np *NodePool) nodeInstanceGroup(n v1.Node, client *http.Client) (string, error) {
	// gets the GKE clusterName name
	clusterName, err := metadata.InstanceAttributeValue("cluster-name")
	if err != nil {
		return "", errors.Wrap(err, "Could not get cluster name from metadata")
	}

	ctnService, err := container.New(client)
	if err != nil {
		return "", errors.Wrap(err, "Could not create container service")
	}

	npName := n.Labels["cloud.google.com/gke-nodepool"]
	nodePool, err := container.NewProjectsZonesClustersNodePoolsService(ctnService).Get(np.project, np.zone, clusterName, npName).Do()
	if err != nil {
		return "", errors.Wrapf(err, "Could not get node pool details from cluster %v, nodepool %v", clusterName, npName)
	}

	// Response: "https://www.googleapis.com/compute/v1/projects/unity-1355/zones/us-west1-b/instanceGroupManagers/gke-soccer-us-west-game-servers-60789234-grp"
	if len(nodePool.InstanceGroupUrls) == 0 {
		return "", errors.Errorf("Could not determine instance group for node: %v", n.Name)
	}
	parts := strings.Split(nodePool.InstanceGroupUrls[0], "/")
	return parts[len(parts)-1], nil
}
