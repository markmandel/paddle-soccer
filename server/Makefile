# Copyright 2016 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# Makefile general project tasks
#

#  __     __         _       _     _
#  \ \   / /_ _ _ __(_) __ _| |__ | | ___ ___
#   \ \ / / _` | '__| |/ _` | '_ \| |/ _ \ __|
#    \ V / (_| | |  | | (_| | |_) | |  __\__ \
#     \_/ \__,_|_|  |_|\__,_|_.__/|_|\___|___/
#

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_path := $(dir $(mkfile_path))

#   _____                    _
#  |_   _|_ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    |_|\__,_|_|  \__, |\___|\__|___/
#                 |___/

# test all the packages
test:
	cd sessions && $(MAKE) test
	cd pkg && $(MAKE) test
	cd matchmaker && $(MAKE) test

# run the visualiser
visualise:
	kubectl proxy --www=/home/$(shell whoami)/gcp-live-k8s-visualizer --www-prefix=/static/ --api-prefix=/api/ --address="0.0.0.0"

# bring up the kubectl proxy if you want to look at it.
kubectl-proxy:
	 kubectl proxy --address=0.0.0.0

# install weavescope, for nicer visualisations
# https://www.weave.works/docs/scope/latest/installing/#k8s
install-weavescope:
	# kubectl apply -f 'https://cloud.weave.works/launch/k8s/weavescope.yaml'
	# make the service public, so it's easier to get to for demos
	wget -O /tmp/weavescope.yaml https://cloud.weave.works/launch/k8s/weavescope.yaml
	sed -i 's/selector:/type: LoadBalancer\n      selector:/' /tmp/weavescope.yaml
	kubectl apply -f /tmp/weavescope.yaml

# delete everything from k8s
clean-k8s:
	kubectl delete deployments --all
	kubectl delete services --all
	kubectl delete statefulset --all
	kubectl delete pods --all
	kubectl delete persistentvolume --all
	kubectl delete persistentvolumeclaim --all
	kubectl delete daemonsets --all

# easy way to pull in all the dependencies the project uses
# handy if godep does something odd.
deps:
	-git clone -b release-2.0 --depth=5 https://github.com/kubernetes/client-go.git $(GOPATH)/src/k8s.io/client-go
	cd $(GOPATH)/src/k8s.io/client-go && godep restore && rm -rf vendor && git add -A && git commit -a -m "Commit"
	go get github.com/stretchr/testify/assert
	go get github.com/gorilla/mux
	go get github.com/garyburd/redigo/redis
	go get github.com/cenkalti/backoff

