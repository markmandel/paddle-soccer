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
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// gameHandler is the handler to post to, such that game match-making can occur
func gameHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	con := s.pool.Get()
	defer con.Close()

	log.Print("[Info][game_route] Match to a game")
	g, err := popOpenGame(con)

	if err != nil {
		if err != errGameNotFound {
			return err
		}

		g = NewGame()
		err := pushOpenGame(con, g)

		if err != nil {
			return err
		}

		// return 201 when pushing into the list
		w.WriteHeader(http.StatusCreated)
	} else {
		// creates the running server, and returns
		// a game with the ip and port populated
		g, err = s.createSessionForGame(con, g)
		if err != nil {
			return err
		}

		// update the record
		err := updateGame(con, g)
		if err != nil {
			return err
		}
	}

	err = json.NewEncoder(w).Encode(g)
	if err != nil {
		log.Printf("[Error][game_route] encoding JSON: %v", err)
		return err
	}

	return nil
}

// getHandler is a handler tp get the details of a game
// that is currently running / waiting for a second person to join
func getHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	con := s.pool.Get()
	defer con.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("[Info][get_route] Retriving game: %v", id)

	g, err := getGame(con, redisGamePrefix+id)

	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(g)
	if err != nil {
		log.Printf("[Error][get_route] Error encoding game to json: %v", err)
	}

	return err
}
