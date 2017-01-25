// Server binary for session management
package main

import (
	"log"
	"os"

	"github.com/markmandel/sessions"
)

const (
	// port to listen on
	port = "PORT"
	// address to listen to redis on
	redisAddress = "REDIS_SERVICE"
)

func main() {
	// get environment variables
	port := os.Getenv(port)
	// default for port
	if port == "" {
		port = "8080"
	}
	redisAddr := os.Getenv(redisAddress)

	log.Print("[Info][Main] Creating server...")
	s := sessions.NewServer(":"+port, redisAddr)
	s.Start()
}
