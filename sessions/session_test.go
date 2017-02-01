package sessions

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func TestServerStoreSession(t *testing.T) {
	s := NewServer("", "")
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
	s := NewServer("", "")
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
