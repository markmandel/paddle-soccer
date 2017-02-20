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

package sessions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// createHandler is a handler for creating a game server session pod
func createHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	log.Printf("[Info][create] creating a game session with image: %v", s.gameServerImage)

	id, err := s.createSessionPod()
	if err != nil {
		return err
	}

	sess := Session{ID: id}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(&sess)
}

// registerHandler is a handler for a new game session to register itself with this
// system, so that we know what port the game server has started up on
func registerHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	var sess Session

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[Error][Register] Reading Body: %v.", err)
		return err
	}

	log.Printf("[Info][Register] Recieved JSON Playload: %v", string(b))

	err = json.Unmarshal(b, &sess)

	if err != nil {
		log.Printf("[Error][Register] Error decoding json: %v. [%v]", err, string(b))
		return err
	}

	log.Printf("[Info][Register] Recieved Session Registration: %#v", sess)

	if err != nil {
		log.Printf("[Error][Register] Error connecting to Kubernetes: %v", err)
		return err
	}

	list, err := s.hostNameAndIP()
	if err != nil {
		return err
	}

	ip, err := s.externalNodeIPofPod(sess, list)
	if err != nil {
		return err
	}

	log.Printf("[Info][Register] Session: IP: %v, Port: %v", ip, sess.Port)
	sess.IP = ip

	return s.storeSession(sess)
}

// getHandler is a handler for retrieving information about a specific
// game server session
func getHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		msg := "No session id provided"
		log.Printf("[Error][get] %v", msg)
		http.Error(w, msg, http.StatusNotFound)
		return nil
	}

	log.Printf("[Info][get] Getting Session: %v", id)

	sess, err := s.getSession(id)

	if err == ErrSessionNotFound {
		http.Error(w, fmt.Sprintf("Could not find session for id: %v", id), http.StatusNotFound)
		return nil
	} else if err != nil {
		log.Printf("[Error][get] Error getting session: %v", err)
		return err
	}

	return json.NewEncoder(w).Encode(sess)
}
