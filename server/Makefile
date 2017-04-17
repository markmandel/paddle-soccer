# Copyright 2017 Google Inc. All Rights Reserved.
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

DEV_IMAGE=markmandel/paddle-soccer-dev
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

# run all the tests in the CI docker container
docker-test:
	docker run --rm \
	-v $(current_path):/go/src/github.com/markmandel/paddle-soccer/server \
	-w="/go/src/github.com/markmandel/paddle-soccer/server" \
	$(DEV_IMAGE) /bin/bash -c "/redis/src/redis-server /redis/redis.conf && make test"

# bring up the kubectl proxy if you want to look at it.
kubectl-proxy:
	 kubectl proxy --address=0.0.0.0

# install weavescope, for nicer visualisations
# https://www.weave.works/docs/scope/latest/installing/#k8s
install-weavescope:
	kubectl apply -f https://cloud.weave.works/launch/k8s/weavescope.yaml?k8s-service-type=LoadBalancer&v=1.3.0
	sleep 5
	kubectl edit deployment weave-scope-app

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
# handy because godep looks in the $GOPATH.
# Uses https://github.com/golang/dep for dependency management
# You may want to also use `govend --prune` for pruning
deps:
	dep ensure k8s.io/client-go@release-2.0
	dep ensure github.com/gorilla/mux
	dep ensure github.com/garyburd/redigo #/redis
	dep ensure github.com/cenkalti/backoff
	dep ensure github.com/pkg/errors
	dep ensure github.com/stretchr/testify #/assert

# Fire up a godoc server
godoc:
	godoc -http=":8080" &
