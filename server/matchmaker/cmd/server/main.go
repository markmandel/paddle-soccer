// Server binary for session management
package main

import (
	"log"
	"os"

	"github.com/markmandel/paddle-soccer/server/matchmaker"
)

const (
	// portEnv is the environment variable to
	// use to find the port to listen on
	portEnv = "PORT"
	// redisAddressEnv is the environment variable to find the
	// address to listen to redis on
	redisServiceEnv = "REDIS_SERVICE"
	// sessionsServiceEnv is the environment variable for the url
	// the session service can be found
	sessionsServiceEnv = "SESSIONS_SERVICE"
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
	s := matchmaker.NewServer(":"+port, os.Getenv(redisServiceEnv), os.Getenv(sessionsServiceEnv))
	s.Start()
}
