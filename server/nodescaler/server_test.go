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

package nodescaler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/rand"
)

func TestServerOptions(t *testing.T) {
	sel := "app=game-server"
	cpuRequest := "0.5"
	s, err := NewServer("", sel, cpuRequest)
	assert.Nil(t, err)
	assert.Equal(t, sel, s.nodeSelector)
	assert.EqualValues(t, 500, s.cpuRequest)

	// property based test for buffer count
	for i := 0; i < 20; i++ {
		bc := rand.Int63nRange(5, 100)
		s, err := NewServer("", sel, cpuRequest, ServerBufferCount(bc))
		assert.Nil(t, err)
		assert.Equal(t, bc, s.bufferCount)
	}

	// property based test for shutdown, also check multiple properties
	for i := 0; i < 20; i++ {
		bc := rand.Int63nRange(5, 100)
		sd := rand.Int63nRange(5, 100)
		s, err := NewServer("", sel, cpuRequest, ServerBufferCount(bc), ServerShutdown(time.Duration(sd)*time.Second))
		assert.Nil(t, err)
		assert.Equal(t, bc, s.bufferCount)
		assert.EqualValues(t, sd, s.shutdown.Seconds())
	}

	// property based test for buffer count
	for i := 0; i < 20; i++ {
		st := rand.Int63nRange(5, 100)
		s, err := NewServer("", sel, cpuRequest, ServerTick(time.Duration(st)*time.Second))
		assert.Nil(t, err)
		assert.EqualValues(t, st, s.tick.Seconds())
	}

	// property based test for min nodes
	for i := 0; i < 20; i++ {
		nc := rand.Int63nRange(1, 100)
		s, err := NewServer("", sel, cpuRequest, ServerMinNodeNumber(nc))
		assert.Nil(t, err)
		assert.Equal(t, nc, s.minNodeNumber)
	}

	// property based test for max nodes
	for i := 0; i < 20; i++ {
		nc := rand.Int63nRange(1, 100)
		s, err := NewServer("", sel, cpuRequest, ServerMaxNodeNumber(nc))
		assert.Nil(t, err)
		assert.Equal(t, nc, s.maxNodeNumber)
	}
}
