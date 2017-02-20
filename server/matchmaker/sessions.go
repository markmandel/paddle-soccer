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

package matchmaker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	// sessionMaxRetries is the number of times to
	// check and see if the session is available.
	sessionMaxRetries = 30
)

// Session represents a game session
type Session struct {
	ID   string `json:"id"`
	Port int    `json:"port,omitempty"`
	IP   string `json:"ip,omitempty"`
}

// createSessionForGame sends a message to the session manager to create
// a new game server session. Waits for the game server IP and Port to become active
// and returns a fully populated Game, with the sessions details
func (s *Server) createSessionForGame(con redis.Conn, g *Game) (*Game, error) {
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

	return g, updateGame(con, g)
}

// getSessionIPAndPort returns a Session with the IP and Port of a running
// game session. Will time out after 30 attempts, with a 1 second wait in between.
func (s *Server) getSessionIPAndPort(sess Session) (Session, error) {
	var body io.ReadCloser

	for i := 0; i <= sessionMaxRetries; i++ {
		req := s.sessionAddr + "/session/" + url.QueryEscape(sess.ID)
		log.Printf("[Info][sessions] Requesting Session Data: %v", req)
		res, err := http.Get(req)
		if err != nil {
			log.Printf("[Error][sessions] Error getting session info: %v", err)
			return sess, err
		}
		if res.StatusCode == http.StatusOK {
			log.Printf("[Info][sessions] Recieved session data, status: %v", res.StatusCode)
			body = res.Body
			break
		}
		err = res.Body.Close()
		if err != nil {
			log.Printf("[Warn][sessions] Could not close body: %v", err)
		}

		log.Printf("[Info][sessions] Session %v data not found, trying again", sess.ID)
		time.Sleep(time.Second)
	}
	defer body.Close()

	if body == nil {
		err := fmt.Errorf("Could not get session %v data", sess.ID)
		log.Printf("[Error][sessions] %v", err)
		return sess, err
	}

	err := json.NewDecoder(body).Decode(&sess)

	if err != nil {
		log.Printf("[Error][sessions] Could not decode json to Session, %v", err)
	}

	return sess, err
}
