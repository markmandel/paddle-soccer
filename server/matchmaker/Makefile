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
# Session service
#

#  __     __         _       _     _
#  \ \   / /_ _ _ __(_) __ _| |__ | | ___ ___
#   \ \ / / _` | '__| |/ _` | '_ \| |/ _ \ __|
#    \ V / (_| | |  | | (_| | |_) | |  __\__ \
#     \_/ \__,_|_|  |_|\__,_|_.__/|_|\___|___/
#

PROJECT?=$(shell gcloud config list project | awk 'FNR ==2 { print $$3 }')
TAG=gcr.io/$(PROJECT)/matchmaker
PACKAGE=$(PACKAGE_ROOT)/matchmaker
SERVER?=localhost:8080

#Directory that this Makefile is in.
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_path := $(dir $(mkfile_path))
src_path = $(current_path)/src

#   _____                    _
#  |_   _|_ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    |_|\__,_|_|  \__, |\___|\__|___/
#                 |___/

include ../common.mk

# deploy the replica set
deploy-rs:
	$(MAKE) FILE=deployment template-apply

deploy-service:
	$(MAKE) FILE=service template-apply

deploy-redis:
	kubectl create -f redis.yaml

clean-k8s:
	-kubectl delete deployment matchmaker
	-kubectl delete statefulset matchmaker-redis
	-kubectl delete service matchmaker
	-kubectl delete service matchmaker-redis

# delete the pods, so they come down again
re-deploy:
	kubectl get pods --selector=role=matchmaker-server | grep matchmaker | awk '{ print $$1 }' | xargs kubectl delete pod

test-join:
	curl -v -H "Content-Type: application/json" -X POST http://$(SERVER)/game