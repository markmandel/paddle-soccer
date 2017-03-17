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
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestConnectRedis(t *testing.T) {
	t.Parallel()
	p := NewPool(":6379")

	err := WaitForConnection(p)
	assert.Nil(t, err, "Could not ping redis: %v", err)
}

func TestNewReadinessCheck(t *testing.T) {
	t.Parallel()
	r := &http.Request{}
	w := httptest.NewRecorder()

	rc := NewReadinessCheck(NewPool(":6379"))
	rc(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	// not a valid port
	w = httptest.NewRecorder()
	rc = NewReadinessCheck(NewPool(":4444"))
	rc(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
