package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectRedis(t *testing.T) {
	p := NewPool(":6379")

	err := WaitForConnection(p)
	assert.Nil(t, err, "Could not ping redis: %v", err)
}
