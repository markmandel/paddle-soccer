package sessions

import (
	"errors"
	"log"

	"github.com/garyburd/redigo/redis"
)

// Session represents a game session
type Session struct {
	ID   string `json:"id" redis:"id"`
	Port int    `json:"port,omitempty" redis:"port"`
	IP   string `json:"ip,omitempty" redis:"ip"`
}

const redisSessionPrefix = "Session:"

// ErrorSessionNotFound is returned when you
// can't find the Session in redis
var ErrorSessionNotFound = errors.New("Could not find the requested session")

// storeSession store the session in redis
func (s *Server) storeSession(sess Session) error {
	log.Print("[Info][session] Storing session in redis")
	con := s.pool.Get()
	defer con.Close()

	key := redisSessionPrefix + sess.ID

	err := con.Send("MULTI")
	if err != nil {
		log.Printf("[Error][session] Could not Send MULTI: %v", err)
	}
	err = con.Send("HMSET", key, "id", sess.ID,
		"port", sess.Port,
		"ip", sess.IP)

	if err != nil {
		log.Printf("[Error][session] Could not Send HMSET: %v", err)
		return err
	}
	err = con.Send("EXPIRE", key, 60*60)
	if err != nil {
		log.Printf("[Error][session] Could not Send EXPIRE: %v", err)
		return err
	}
	_, err = con.Do("EXEC")

	if err != nil {
		log.Printf("[Error][session] Could not save session to redis: %v", err)
	}

	return err
}

func (s *Server) getSession(id string) (Session, error) {
	con := s.pool.Get()
	defer con.Close()

	key := redisSessionPrefix + id

	log.Printf("[Info][session_type] Getting data for key: %v", key)

	var result Session
	values, err := redis.Values(con.Do("HGETALL", key))

	if err != nil {
		log.Printf("[Error][session] Error getting hash for key %v. %v", key, err)
		return result, err
	}

	if len(values) == 0 {
		log.Printf("[Error][session] Could not find record for key: %v", key)
		return result, ErrorSessionNotFound
	}

	err = redis.ScanStruct(values, &result)

	if err != nil {
		log.Printf("[Error][session] Error Scanning Struct %v", err)
	}

	return result, err
}
