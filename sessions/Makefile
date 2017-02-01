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


#  __     __         _       _     _
#  \ \   / /_ _ _ __(_) __ _| |__ | | ___ ___
#   \ \ / / _` | '__| |/ _` | '_ \| |/ _ \ __|
#    \ V / (_| | |  | | (_| | |_) | |  __\__ \
#     \_/ \__,_|_|  |_|\__,_|_.__/|_|\___|___/
#

PROJECT?=$(shell gcloud config list project | awk 'FNR ==2 { print $$3 }')
TAG=gcr.io/$(PROJECT)/sessions
PACKAGE_ROOT=github.com/markmandel/sessions
SERVER?=localhost:8080

#Directory that this Makefile is in.
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_path := $(dir $(mkfile_path))
src_path = $(current_path)/src

# build the docker image
build: build-static-server
		docker build --tag=$(TAG) $(current_path)

# deploy the replica set
deploy-rs:
	$(MAKE) FILE=deployment template-apply

deploy-service:
	$(MAKE) FILE=service template-apply

deploy-redis:
	kubectl create -f redis.yaml

clean-k8s:
	kubectl delete deployment sessions

# give me a shell
shell:
	docker run -it --entrypoint=/bin/sh $(TAG)

# run it!
run:
	docker run --rm -p 8080:8080 $(TAG)

# clean all the things
clean: clean-server
	docker rmi $(TAG)

# push to gcr.io
push:
	gcloud docker -- push $(TAG)

# delete the pods, so they come down again
re-deploy:
	kubectl get pods --selector=role=sessions-server | grep sessions | awk '{ print $$1 }' | xargs kubectl delete pod

# pull from gcr.io
pull:
	gcloud docker pull $(TAG)

# Build a statically compiled binary
# https://github.com/constabulary/gb/issues/328
build-static-server: bin-dir
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w -extld ld -extldflags -static' -o $(current_path)/bin/server $(PACKAGE_ROOT)/cmd/server

# Build a go binary. For testing, since we need a statically compiled binary
build-server: bin-dir
	go build -o $(current_path)/bin/server $(PACKAGE_ROOT)/cmd/server

# Make sure the code passes tests
check-code:
	-goimports -w $(current_path) && \
	go vet $(PACKAGE_ROOT) && \
	golint $(PACKAGE_ROOT) && \
	errcheck $(PACKAGE_ROOT) && \
	echo "...Passed"

test: check-code
	goapp test $(PACKAGE_ROOT)

# Fire up a godoc server
godoc:
	godoc -http=":8080" &

deps:
	-git clone -b release-2.0 --depth=5 https://github.com/kubernetes/client-go.git $(GOPATH)/src/k8s.io/client-go
	cd $(GOPATH)/src/k8s.io/client-go && godep restore && rm -rf vendor && git add -A && git commit -a -m "Commit"
	go get github.com/stretchr/testify/assert
	go get github.com/gorilla/mux
	go get github.com/garyburd/redigo/redis
	go get github.com/cenkalti/backoff

test-register:
	curl -v -H "Content-Type: application/json" -X POST -d '{"id":"game-1234","port":7878}' http://$(SERVER)/register

test-create:
	curl -v -H "Content-Type: application/json" -X POST http://$(SERVER)/session

bin-dir:
	-mkdir -p $(current_path)/bin

template-apply:
		@cp $(FILE).yaml /tmp/
		@sed -i 's/$${PROJECT}/$(PROJECT)/g' $(FILE).yaml
		kubectl apply -f $(FILE).yaml --record
		@cp $(FILE).yaml /tmp/$(FILE).deployed.yaml
		@cp /tmp/$(FILE).yaml .

