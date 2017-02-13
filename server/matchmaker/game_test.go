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

package matchmaker

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	g := NewGame()
	assert.NotNil(t, g.ID)

	assert.Equal(t, g.Status, gameStatusOpen)
}

func TestPushOpenGame(t *testing.T) {
	s := NewServer("", "", "")
	defer s.pool.Close()
	con := s.pool.Get()

	g := NewGame()
	err := pushOpenGame(con, g)
	assert.Nil(t, err)

	key := g.Key()
	list, err := redis.Strings(con.Do("LRANGE", redisOpenGameListKey, 0, -1))
	assert.Nil(t, err)

	found := false
	for _, i := range list {
		if i == key {
			found = true
			break
		}
	}

	assert.True(t, found, "Key was not in the list %v", list)

	exist, err := redis.Bool(con.Do("EXISTS", key))
	assert.Nil(t, err)
	assert.True(t, exist)

	tty, err := redis.Int(con.Do("TTL", key))
	assert.Nil(t, err)
	assert.Equal(t, 3600, tty)
}

func TestPopOpenGame(t *testing.T) {
	s := NewServer("", "", "")
	defer s.pool.Close()
	con := s.pool.Get()

	_, err := con.Do("FLUSHALL")
	assert.Nil(t, err)

	g, err := popOpenGame(con)
	assert.Nil(t, g)
	assert.Equal(t, errGameNotFound, err)

	g = NewGame()
	err = pushOpenGame(con, g)
	assert.Nil(t, err)

	reget, err := popOpenGame(con)
	assert.Nil(t, err)
	assert.Equal(t, g, reget)
}

func TestGetAndUpdateGame(t *testing.T) {
	s := NewServer("", "", "")
	defer s.pool.Close()
	con := s.pool.Get()

	_, err := con.Do("FLUSHALL")
	assert.Nil(t, err)

	g := NewGame()
	err = pushOpenGame(con, g)
	assert.Nil(t, err)

	reget, err := getGame(con, g.Key())
	assert.Nil(t, err)
	assert.Equal(t, g, reget)

	g.Status = gameStatusClosed
	g.Port = 8080
	g.IP = "1.2.3.4"

	err = updateGame(con, g)
	assert.Nil(t, err)

	reget, err = getGame(con, g.Key())
	assert.Nil(t, err)
	assert.Equal(t, g, reget)
}
