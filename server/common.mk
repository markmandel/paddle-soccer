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
# Common makefile targets and functions
#

#  __     __         _       _     _
#  \ \   / /_ _ _ __(_) __ _| |__ | | ___ ___
#   \ \ / / _` | '__| |/ _` | '_ \| |/ _ \ __|
#    \ V / (_| | |  | | (_| | |_) | |  __\__ \
#     \_/ \__,_|_|  |_|\__,_|_.__/|_|\___|___/
#



#   _____                    _
#  |_   _|_ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    |_|\__,_|_|  \__, |\___|\__|___/
#                 |___/

# build the docker image
build: build-static-server
	docker build --tag=$(TAG) $(current_path)

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
	find $(current_path) -type f -name '*.go' -not -path '$(current_path)vendor/*' | xargs goimports -w
	go vet $(PACKAGE_ROOT)
	golint $(PACKAGE_ROOT)
	-errcheck $(PACKAGE_ROOT) | grep -v "defer"
	@echo "...Complete"

test: check-code
	go test $(PACKAGE_ROOT)

bin-dir:
	-mkdir -p $(current_path)/bin

template-apply:
		@cp $(FILE).yaml /tmp/
		@sed -i 's/$${PROJECT}/$(PROJECT)/g' $(FILE).yaml
		kubectl apply -f $(FILE).yaml --record
		@cp $(FILE).yaml /tmp/$(FILE).deployed.yaml
		@cp /tmp/$(FILE).yaml .
