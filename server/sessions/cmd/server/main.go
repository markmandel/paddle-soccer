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
	s.Start()
}
