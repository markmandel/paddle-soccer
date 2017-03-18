package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeserialiseEnvMap(t *testing.T) {
	m := deserialiseEnvMap("")
	assert.Equal(t, map[string]string{}, m)

	m = deserialiseEnvMap("foo:bar")
	assert.Equal(t, map[string]string{"foo": "bar"}, m)

	m = deserialiseEnvMap("foo:bar,frog:scorpion")
	assert.Equal(t, map[string]string{"foo": "bar", "frog": "scorpion"}, m)
}
