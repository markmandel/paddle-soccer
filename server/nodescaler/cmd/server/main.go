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

	"time"

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

	// tickEnv is the environment variable to tell
	// the scaler the duration between each check
	tickerEnv = "TICK"

	// shutDownNodeEnv is the environment variable to tell
	// the scaler how long after a node is cordoned, should it be shut down
	// (once it is empty)
	shutDownNodeEnv = "SHUTDOWN_NODE"

	// minNodeEnv is the environment variable that controls what the
	// minimum number of nodes in the cluster have to be at any given
	// point and time
	minNodeEnv = "MIN_NODE"

	// maxNodeEnv is the environment variable that controls what the
	// maximum number of nodes in the cluster have can be at any given
	// point and time
	maxNodeEnv = "MAX_NODE"
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

	var opts []nodescaler.Option

	if bc := os.Getenv(bufferCountEnv); bc != "" {
		bufferCount, err := strconv.ParseInt(bc, 10, 64)
		if err != nil {
			log.Fatalf("[Error][Main] Error decoding %v value of %v, %v", bufferCountEnv, bc, err)
		}
		opts = append(opts, nodescaler.ServerBufferCount(bufferCount))
	}

	if t := os.Getenv(tickerEnv); t != "" {
		tick, err := time.ParseDuration(t)
		if err != nil {
			log.Fatalf("[Error][Main] Error parsing %v value of %v, %v", tickerEnv, t, err)
		}
		opts = append(opts, nodescaler.ServerTick(tick))
	}

	if sd := os.Getenv(shutDownNodeEnv); sd != "" {
		shutDown, err := time.ParseDuration(sd)
		if err != nil {
			log.Fatalf("[Error][Main] Error decoding %v value of %v, %v", shutDownNodeEnv, sd, err)
		}
		opts = append(opts, nodescaler.ServerShutdown(shutDown))
	}

	if min := os.Getenv(minNodeEnv); min != "" {
		count, err := strconv.ParseInt(min, 10, 64)
		if err != nil {
			log.Fatalf("[Error][Main] Error decoding %v value of %v, %v", minNodeEnv, min, err)
		}
		opts = append(opts, nodescaler.ServerMinNodeNumber(int64(count)))
	}

	if max := os.Getenv(maxNodeEnv); max != "" {
		count, err := strconv.ParseInt(max, 10, 64)
		if err != nil {
			log.Fatalf("[Error][Main] Error decoding %v value of %v, %v", maxNodeEnv, max, err)
		}
		opts = append(opts, nodescaler.ServerMaxNodeNumber(int64(count)))
	}

	s, err := nodescaler.NewServer(":"+port, os.Getenv(nodeSelectorEnv), os.Getenv(cpuRequestEnv), opts...)
	if err != nil {
		log.Fatalf("[Error][Main] %+v", err)
	}

	if err := s.Start(); err != nil {
		log.Fatalf("[Error][Main] %+v", err)
	}
}
