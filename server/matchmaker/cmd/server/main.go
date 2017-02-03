// Server binary for session management
package main

import (
	"log"
	"os"

	"github.com/markmandel/paddle-soccer/matchmaker"
)

const (
	// port to listen on
	portEnv = "PORT"
	// address to listen to redis on
	redisServiceEnv = "REDIS_SERVICE"
	// sessions address
	sessionsServiceEnv = "SESSIONS_SERVICE"
)

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
