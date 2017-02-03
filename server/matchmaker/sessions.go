package matchmaker

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/garyburd/redigo/redis"
)

// Session represents a game session
type Session struct {
	ID   string `json:"id"`
	Port int    `json:"port,omitempty"`
	IP   string `json:"ip,omitempty"`
}

func (s *Server) createSession(con redis.Conn, g *Game) (*Game, error) {
	path := s.sessionAddr + "/session"
	log.Printf("[Info][sessions] Creating a new session: %v", path)
	r, err := http.Post(path, "application/json", nil)

	if err != nil {
		log.Printf("[Error][sessions] Error calling /session: %v", err)
		return g, err
	}
	defer r.Body.Close()

	sess := Session{}
	err = json.NewDecoder(r.Body).Decode(&sess)
	if err != nil {
		log.Printf("[Error][sessions] Error: %v", err)
		return g, err
	}

	log.Printf("[Info][sessions] Created Session: %#v", sess)
	g.SessionID = sess.ID

	sess, err = s.getSessionIPAndPort(sess)
	g.Port = sess.Port
	g.IP = sess.IP
	g.Status = gameStatusClosed

	err = updateGame(con, g)
	if err != nil {
		return g, err
	}

	return g, nil
}

func (s *Server) getSessionIPAndPort(sess Session) (Session, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Minute

	var body io.ReadCloser

	err := backoff.Retry(func() error {
		req := s.sessionAddr + "/session/" + url.QueryEscape(sess.ID)
		log.Printf("[Info][sessions] Requesting Session Data: %v", req)
		res, err := http.Get(req)
		if err != nil {
			log.Printf("[Error][sessions] Error getting session info: %v", err)
			return err
		}
		if res.StatusCode == http.StatusNotFound {
			log.Printf("[Info][sessions] Session %v data not found, trying again", sess.ID)
			return errors.New("Not found. Try again")
		}

		log.Printf("[Info][sessions] Recieved session data, status: %v", res.StatusCode)

		body = res.Body
		return nil
	}, bo)
	defer body.Close()

	if err != nil {
		log.Printf("[Error][sessions] Could not get session %v data, %v", sess.ID, err)
	}

	err = json.NewDecoder(body).Decode(&sess)

	if err != nil {
		log.Printf("[Error][sessions] Could not decode json to Session, %v", err)
	}

	return sess, err
}
