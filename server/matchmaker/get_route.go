package matchmaker

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	con := s.pool.Get()
	defer con.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("[Info][get_route] Retriving game: %v", id)

	g, err := getGame(con, redisGamePrefix+id)

	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(g)
	if err != nil {
		log.Printf("[Error][get_route] Error encoding game to json: %v", err)
	}

	return err
}
