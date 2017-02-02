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
local_gopath := $(current_path)/go
project_path := $(local_gopath)/src/github.com/markmandel/paddle-soccer

#   _____                    _
#  |_   _|_ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    |_|\__,_|_|  \__, |\___|\__|___/
#                 |___/

# Create create a local gopath directory, with the go
# projects in the right places.
# This is already set to be ignored by git
create-local-gopath:
	mkdir -p $(local_gopath)/src/github.com/markmandel
	ln -sr $(current_path) $(project_path)

# cleans the local gopath
clean-local-gopath:
	rm -rf $(local_gopath)

# test all the packages
test:
	cd sessions && $(MAKE) test
	cd pkg && $(MAKE) test

# Fire up a godoc server
godoc:
	godoc -http=":8080" &

# easy way to pull in all the dependencies the project uses
# handy if godep does something odd.
deps:
	-git clone -b release-2.0 --depth=5 https://github.com/kubernetes/client-go.git $(GOPATH)/src/k8s.io/client-go
	cd $(GOPATH)/src/k8s.io/client-go && godep restore && rm -rf vendor && git add -A && git commit -a -m "Commit"
	go get github.com/stretchr/testify/assert
	go get github.com/gorilla/mux
	go get github.com/garyburd/redigo/redis
	go get github.com/cenkalti/backoff

