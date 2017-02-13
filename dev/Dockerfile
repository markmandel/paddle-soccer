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

FROM golang:1.7

RUN apt-get update && \
    apt-get -y install wget curl entr nano less unzip python zsh git make openssh-server libapparmor1 libltdl7

WORKDIR /

#sshd setup - https://docs.docker.com/examples/running_ssh_service/
RUN mkdir /var/run/sshd
RUN echo 'root:pw' | chpasswd
RUN sed -i 's/PermitRootLogin without-password/PermitRootLogin yes/' /etc/ssh/sshd_config
# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd
ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile
EXPOSE 22

#install cloud sdk
RUN wget -q https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.zip && unzip -q google-cloud-sdk.zip && rm google-cloud-sdk.zip
RUN /google-cloud-sdk/install.sh --usage-reporting=true --path-update=true --bash-completion=true --rc-path=/root/.bashrc
ENV PATH /google-cloud-sdk/bin:$PATH

RUN gcloud components update kubectl

#redis
RUN mkdir redis && cd redis && wget -O redis.tar.gz -q http://download.redis.io/releases/redis-3.2.6.tar.gz && \
    tar xzf redis.tar.gz --strip-components=1 && rm redis.tar.gz && make
RUN sed -i "s/daemonize no/daemonize yes/" /redis/redis.conf
RUN sed -i 's|logfile ""|logfile "/redis/redis.log"|' /redis/redis.conf
RUN sed -i 's|bind 127.0.0.1|bind 0.0.0.0|' /redis/redis.conf

#oh-my-zsh, because how do we live without it?
RUN git clone https://github.com/robbyrussell/oh-my-zsh.git

#tools for working with go.
RUN go get github.com/golang/lint/golint
RUN go get golang.org/x/tools/cmd/goimports
RUN go get github.com/kisielk/errcheck
RUN go get github.com/tools/godep

# copy these tools into Go install, so we can mount on /go
# and when we mount the go sdk they are available
# externally
RUN mv /go/bin/* /usr/local/go/bin/
RUN rm -r /go/src
RUN rm -r /go/pkg

WORKDIR /root

#install kubectl visualiser
RUN git clone --depth=5 -b "kubernetes-1.2" https://github.com/saturnism/gcp-live-k8s-visualizer.git
# K8s proxy
EXPOSE 8001

ADD startup.sh startup.sh
RUN chmod +x startup.sh

EXPOSE 8080