#!/usr/bin/env sh

# Copyright 2016 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

groupadd --gid $HOST_GID $HOST_USER
useradd $HOST_USER --home /home/$HOST_USER --gid $HOST_GID --uid $HOST_UID --shell /usr/bin/zsh
echo "$HOST_USER:pw" | chpasswd

chown -R $HOST_USER:$HOST_USER /home/$HOST_USER
chown -R $HOST_USER:$HOST_USER /oh-my-zsh
chown -R $HOST_USER:$HOST_USER /google-cloud-sdk
chown -R $HOST_USER:$HOST_USER /go

#allow docker passthrough
groupadd --gid $DOCKER_GID docker
usermod -a -G docker $HOST_USER

#link up go src, so it works
ln -s /home/$HOST_USER/project/server/go/src /go

#start redis
/redis/src/redis-server /redis/redis.conf

/usr/sbin/sshd
su $HOST_USER
