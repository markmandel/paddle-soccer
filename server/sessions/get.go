package sessions

import (
	"encoding/json"
	"log"
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
)

// getHandler registration of a new game session
func getHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		msg := "No session id provided"
		log.Printf("[Error][get]" + msg)
		http.Error(w, msg, http.StatusNotFound)
		return nil
	}

	log.Printf("[Info][get] Getting Session: %v", id)

	sess, err := s.getSession(id)

	if err != nil {
		log.Printf("[Error][get] Error getting session: %v", err)

		if err == ErrorSessionNotFound {
			http.Error(w, fmt.Sprintf("Could not find session for id: %v", id), http.StatusNotFound)
			return nil
		}

		return err
	}

	return json.NewEncoder(w).Encode(sess)
}
