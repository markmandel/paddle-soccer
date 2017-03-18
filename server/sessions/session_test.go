// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sessions

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func TestServerStoreSession(t *testing.T) {
	s := NewServer("", "", "", nil)
	con := s.pool.Get()
	defer con.Close()

	id := uuid.New()
	sess := Session{ID: id, Port: 8080, IP: "0.0.0.0"}

	err := s.storeSession(sess)
	assert.Nil(t, err)

	key := redisSessionPrefix + sess.ID

	exist, err := redis.Bool(con.Do("EXISTS", key))
	assert.Nil(t, err)
	assert.True(t, exist)

	tty, err := redis.Int(con.Do("TTL", key))
	assert.Nil(t, err)
	assert.Equal(t, 3600, tty)
}

func TestServerGetSession(t *testing.T) {
	s := NewServer("", "", "", nil)
	con := s.pool.Get()
	defer con.Close()

	id := uuid.New()
	sess := Session{ID: id, Port: 8080, IP: "0.0.0.0"}

	err := s.storeSession(sess)
	assert.Nil(t, err)

	reget, err := s.getSession(id)
	assert.Nil(t, err)
	assert.Equal(t, sess, reget)
}
