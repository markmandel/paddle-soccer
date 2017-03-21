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

	"strconv"

	"github.com/markmandel/paddle-soccer/server/nodescaler"
)

const (
	// portEnv is the environment variable to
	// use to find the port to listen on
	portEnv = "PORT"

	// nodeSelectorEnv is the environment variable
	// to use to specify (via K8s selector) which Kubernetes nodes
	// that this scaler should manages
	nodeSelectorEnv = "NODE_SELECTOR"

	// cpuRequestEnv is the environment variable to tell the
	// scaler what cpu request size has been selected for each
	// game server instance
	cpuRequestEnv = "CPU_REQUEST"

	// bufferCountEnv is the environment variable to tell
	// the scaler how many game servers you want a buffer of
	// to ensure there is space for new servers.
	bufferCountEnv = "BUFFER_COUNT"
)

// main function for starting the server
func main() {
	// get environment variables
	port := os.Getenv(portEnv)
	// default for portEnv
	if port == "" {
		port = "8080"
	}

	log.Print("[Info][Main] Creating server...")

	bufferCount, err := strconv.Atoi(os.Getenv(bufferCountEnv))
	if err != nil {
		log.Fatalf("[Error][Main] Error decoding %v value of %v, %v", bufferCountEnv, os.Getenv(bufferCountEnv), err)
	}

	s, err := nodescaler.NewServer(":"+port, os.Getenv(nodeSelectorEnv), os.Getenv(cpuRequestEnv), int64(bufferCount))
	if err != nil {
		log.Fatalf("[Error][Main] %+v", err)
	}

	if err := s.Start(); err != nil {
		log.Fatalf("[Error][Main] %+v", err)
	}
}
