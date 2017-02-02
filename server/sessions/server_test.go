package sessions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRedis(t *testing.T) {
	s := NewServer(":8080", "")
	assert.NotNil(t, s)
	err := s.pingRedis()
	assert.Nil(t, err, "Could not ping redis: %v", err)
}
