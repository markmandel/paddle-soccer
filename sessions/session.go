package sessions

import "log"

// Session represents a game session
type Session struct {
	ID   string `json:"id",redis:"id"`
	Port int    `json:"port",redis:"port"`
}

const redisSessionPrefix = "Session:"

// storeSession store the session in redis
func (s *Server) storeSession(sess Session) error {
	log.Print("[Info][Session] Storing session in redis")
	con := s.pool.Get()
	defer con.Close()

	key := redisSessionPrefix + sess.ID

	err := con.Send("MULTI")
	if err != nil {
		log.Printf("[Error][session] Could not Send MULTI: %v", err)
	}
	err = con.Send("HMSET", key, "id", sess.ID, "port", sess.Port)
	if err != nil {
		log.Printf("[Error][session] Could not Send HMSET: %v", err)
	}
	err = con.Send("EXPIRE", key, 60*60)
	if err != nil {
		log.Printf("[Error][session] Could not Send EXPIRE: %v", err)
	}
	_, err = con.Do("EXEC")

	if err != nil {
		log.Printf("[Error][session] Could not save session to redis: %v", err)
	}

	return err
}
