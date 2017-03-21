package sessions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerStart(t *testing.T) {
	_, err := NewServer("", "", "", nil, "0.5")
	assert.Nil(t, err, "0.5 should be fine")

	_, err = NewServer("", "", "", nil, "500m")
	assert.Nil(t, err, "500m should be fine")

	_, err = NewServer("", "", "", nil, "XXXXXX")
	assert.NotNil(t, err, "How is XXXXX valid?")
}
