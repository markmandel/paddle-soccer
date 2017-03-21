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

	"strings"

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

	// gameNodeSelectorEnv is the environment variable
	// that specifies the node selector map.
	// Each entry is separated by , and each key, value
	// is seperate by :
	gameNodeSelectorEnv = "GAME_NODE_SELECTOR"

	// cpuLimitEnv is the environment variable to specify
	// the cpu limits for each game server - see
	// https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#meaning-of-cpu
	// for more details
	cpuLimitEnv = "CPU_LIMIT"
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
	s, err := sessions.NewServer(":"+port, os.Getenv(redisAddressEnv),
		os.Getenv(gameServerImageEnv), deserialiseEnvMap(os.Getenv(gameNodeSelectorEnv)),
		os.Getenv(cpuLimitEnv))

	if err != nil {
		log.Fatalf("[Error][Main] %+v", err)
	}

	if err := s.Start(); err != nil {
		log.Fatalf("[Error][Main] %+v", err)
	}
}

// deserialiseEnvMap takes a string of values, delineated
// by commas for each value, and : to split keys and values
func deserialiseEnvMap(s string) map[string]string {
	result := map[string]string{}

	for _, item := range strings.Split(s, ",") {
		entry := strings.Split(item, ":")
		if len(entry) >= 2 {
			result[entry[0]] = entry[1]
		}
	}

	return result
}
