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
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/markmandel/paddle-soccer/server/pkg"
)

// Version is the current api version number
const Version string = "matchmaker:0.1"

// Server is the http server instance
type Server struct {
	srv         *http.Server
	pool        *redis.Pool
	sessionAddr string
}

// Handler is the extended http.HandleFunc to provide context for this application
type Handler func(*Server, http.ResponseWriter, *http.Request) error

// NewServer returns the HTTP Server instance
func NewServer(hostAddr, redisAddr string, sessionAddr string) *Server {
	if redisAddr == "" {
		redisAddr = ":6379"
	}

	log.Printf("[Info][Server] Starting server version %v on port %v", Version, hostAddr)
	log.Printf("[Info][Server] Connecting to Redis at %v", redisAddr)
	log.Printf("[Info][Server] Connecting to Sessions at %v", sessionAddr)

	s := &Server{pool: pkg.NewPool(redisAddr), sessionAddr: sessionAddr}

	r := s.createRoutes()

	s.srv = &http.Server{
		Handler: r,
		Addr:    hostAddr,
	}

	return s
}

// Start starts the HTTP server on the given port
func (s *Server) Start() {
	err := pkg.WaitForConnection(s.pool)
	if err != nil {
		log.Fatalf("[Error][Server] Could not connect to redis: %v", err)
	}

	log.Fatalf("[Error][Server] Error starting server: %v", s.srv.ListenAndServe())
}

// createRoutes creates the http routes for this application
func (s *Server) createRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/game", s.standardHandler(gameHandler)).Methods("POST")
	r.HandleFunc("/game/{id}", s.standardHandler(getHandler)).Methods("GET")

	return r
}

// standardHandler returns a http.HandleFunc
// wrapped in standard middleware
func (s *Server) standardHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(s, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
