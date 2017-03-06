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

// Server binary for session management
package main

import (
	"log"
	"os"

	"github.com/markmandel/paddle-soccer/server/sessions"
)

const (
	// portEnv is the environment variable to
	// use to find the port to listen on
	portEnv = "PORT"
	// redisAddressEnv is the environment variable to find the
	// address to listen to redis on
	redisAddressEnv = "REDIS_SERVICE"

	// gameServerImageEnv is the environment variable
	// to set the image that the game server should use
	// when starting up games
	gameServerImageEnv = "GAME_SERVER_IMAGE"
)

// main starts the sessions server
func main() {
	// get environment variables
	port := os.Getenv(portEnv)
	// default for port
	if port == "" {
		port = "8080"
	}
	log.Print("[Info][Main] Creating server...")
	s := sessions.NewServer(":"+port, os.Getenv(redisAddressEnv), os.Getenv(gameServerImageEnv))
	if err := s.Start(); err != nil {
		log.Fatalf("[Error][Main] %+v", err)
	}
}
