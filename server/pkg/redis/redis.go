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

package redis

import (
	"log"
	"time"

	"net/http"

	"github.com/cenkalti/backoff"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

// WaitForConnection pings redis with an exponential backoff,
// to wait until we are connected to redis.
func WaitForConnection(pool *redis.Pool) error {

	return backoff.Retry(func() error {
		con := pool.Get()
		defer con.Close()

		_, err := con.Do("PING")
		if err != nil {
			log.Printf("[Warn][Redis] Retrying connection to redis. %v", err)
		} else {
			log.Print("[Info][Redis] Connected.")
		}

		return errors.Wrap(err, "Could not connect to Redis")
	}, backoff.NewExponentialBackOff())
}

// NewPool returns a new redis pool with the standard configuration
func NewPool(address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 4 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return errors.Wrap(err, "Error with PING on TestOnBorrow")
		},
	}
}

// NewReadinessCheck creates a http rediness check
// that tests that a connection to Redis can be made
func NewReadinessCheck(pool *redis.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		con := pool.Get()
		defer con.Close()

		_, err := con.Do("PING")
		if err != nil {
			log.Printf("[Warn][Redis] Readiness check failed. %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
