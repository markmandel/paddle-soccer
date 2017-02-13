# Copyright 2017 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#

#
# Tools for building the server as a Docker image
#


#  __     __         _       _     _
#  \ \   / /_ _ _ __(_) __ _| |__ | | ___ ___
#   \ \ / / _` | '__| |/ _` | '_ \| |/ _ \ __|
#    \ V / (_| | |  | | (_| | |_) | |  __\__ \
#     \_/ \__,_|_|  |_|\__,_|_.__/|_|\___|___/
#

PROJECT?=$(shell gcloud config list project | awk 'FNR ==2 { print $$3 }')
TAG=soccer-server
RUN_NAME=soccer-server
VERSION=0.1
REGISTRY_TAG=gcr.io/$(PROJECT)/$(TAG):$(VERSION)
UNITY_EDITOR=/opt/Unity/Editor/Unity

#Directory that this Makefile is in.
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_path := $(dir $(mkfile_path))
unity_src := $(current_path)../../unity
unity_build := $(unity_src)/Builds/

#   _____                    _
#  |_   _|_ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    |_|\__,_|_|  \__, |\___|\__|___/
#                 |___/

# build the docker image from a clean build
build: build-unity build-docker

# build just the docker image
build-docker: build-server-tar
	docker build --tag=$(TAG) $(current_path)

# tar up the server build
build-server-tar:
	-rm Server.tar.gz
	cd $(unity_build)/Linux && \
	tar -czf Server.tar.gz Server.x86_64 Server_Data && \
	mv Server.tar.gz $(current_path)

clean-docker:
	docker rmi $(TAG)
	docker rmi $(REGISTRY_TAG)

# run it!
run:
	docker run --rm --name $(RUN_NAME) -e "MIN_PORT=7777" -e "MAX_PORT=7777" -p 7777:7777/udp $(TAG)

stop:
	docker stop $(RUN_NAME)

attach:
	docker exec -it $(RUN_NAME) bash

push:
	docker tag $(TAG) $(REGISTRY_TAG)
	gcloud docker -- push $(REGISTRY_TAG)

# interactive shell
shell:
	docker run --rm -it --entrypoint=bash --user=root $(TAG)

# run the native server
run-server:
	$(unity_build)/Linux/Server.x86_64 -logFile /dev/stdout

# build it from unity
build-unity:
	$(UNITY_EDITOR) -batchmode -projectPath "$(unity_src)" -executeMethod Plugins.Editor.Build.BuildGameAndServer -quit -logFile .unity-build.log

# clean the builds directory
clean-unity:
	trash $(unity_build)/*

# run the tests
test:
	$(UNITY_EDITOR) -batchmode -projectPath "$(unity_src)" -runEditorTests -editorTestsResultFile "/tmp/soccer-multiplayer"