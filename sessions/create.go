package sessions

import (
	"log"
	"net/http"

	"encoding/json"
	"os"
)

const gameServerImageEnv = "GAME_SERVER_IMAGE"

// createHandler registration of a new game session
func createHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	image := os.Getenv(gameServerImageEnv)

	log.Printf("[Info][create] creating a game session with image: %v", image)

	id, err := s.createSessionPod(image)
	if err != nil {
		return err
	}

	sess := Session{ID: id}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(&sess)
}
