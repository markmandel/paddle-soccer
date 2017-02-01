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
# Makefile for the shell development environment
#


#  __     __         _       _     _
#  \ \   / /_ _ _ __(_) __ _| |__ | | ___ ___
#   \ \ / / _` | '__| |/ _` | '_ \| |/ _ \ __|
#    \ V / (_| | |  | | (_| | |_) | |  __\__ \
#     \_/ \__,_|_|  |_|\__,_|_.__/|_|\___|___/
#

PROJECT=$(shell gcloud config list project | awk 'FNR ==2 { print $$3 }')
TAG=multiplayer-soccer/dev
REGISTRY_TAG=gcr.io/$(PROJECT)/$(TAG)
NAME=multiplayer-soccer-shell
DOCKER_GID=$(word 3,$(subst :, ,$(shell getent group docker)))

#Directory that this Makefile is in.
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_path := $(dir $(mkfile_path))
project_path := $(realpath $(current_path)/..)

#   _____                    _
#  |_   _|_ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    |_|\__,_|_|  \__, |\___|\__|___/
#                 |___/

# build the docker dev image
build:
	docker build --tag=$(TAG) $(current_path)

push:
	gcloud docker push $(TAG)

pull:
	gcloud docker pull $(TAG)

#clean up the shell
clean:
	docker rmi $(TAG)
	docker rmi $(REGISTRY_TAG)

# Start a dev shell
shell:
	mkdir -p ~/.config/gcloud
	docker run --rm \
		--name=$(NAME) \
		-P=true \
		-p 6379:6379 \
		-e TERM \
		-e HOST_GID=`id -g` \
		-e HOST_UID=`id -u` \
		-e HOST_USER=$(USER) \
		-e DOCKER_GID=$(DOCKER_GID) \
		-v /usr/bin/docker:/usr/bin/docker \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v ~/.gitconfig:/home/$(USER)/.gitconfig \
		-v ~/.ssh:/home/$(USER)/.ssh \
		-v ~/.config/gcloud:/home/$(USER)/.config/gcloud \
		-v ~/.appcfg_oauth2_tokens:/home/$(USER)/.appcfg_oauth2_tokens \
		-v $(project_path)/.kube:/home/$(USER)/.kube \
		-v $(project_path)/dev/zshrc:/home/$(USER)/.zshrc \
		-v $(project_path):/home/$(USER)/project \
		-w /home/$(USER)/project \
		-it $(TAG) /root/startup.sh

# attach to a running dev shell
shell-attach:
	docker exec -it --user=$(USER) $(NAME) zsh

# mount sdks
mount:
	-mkdir -p /tmp/$(NAME)/go
	sshfs $(USER)@0.0.0.0:/usr/local/go /tmp/$(NAME)/go -p $(call getPort,22) -o follow_symlinks

# remove known_hosts of type [0.0.0.0]
clean-mounts:
	sed -i 's/^\[0.0.0.0\].*$$//g' ~/.ssh/known_hosts

#   _____                 _   _
#  |  ___|   _ _ __   ___| |_(_) ___  _ __  ___
#  | |_ | | | | '_ \ / __| __| |/ _ \| '_ \/ __|
#  |  _|| |_| | | | | (__| |_| | (_) | | | \__ \
#  |_|   \__,_|_| |_|\___|\__|_|\___/|_| |_|___/
#

# get the mapped docker host port
getPort = $(word 2,$(subst :, ,$(shell docker port $(NAME) $(1))))
