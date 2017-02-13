package pkg

import (
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/garyburd/redigo/redis"
)

// WaitForConnection pings redis with an exponential backoff,
// to wait until we are connected to redis.
func WaitForConnection(pool *redis.Pool) error {

	return backoff.Retry(func() error {
		con := pool.Get()
		defer con.Close()

		_, err := con.Do("PING")
		if err != nil {
			log.Printf("[Warn][Redis] Could not connect to Redis. %v", err)
		} else {
			log.Print("[Info][Redis] Connected.")
		}

		return err
	}, backoff.NewExponentialBackOff())
}

// NewPool returns a new redis pool with the standard configuration
func NewPool(address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
