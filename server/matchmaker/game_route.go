package matchmaker

import (
	"encoding/json"
	"log"
	"net/http"
)

// GameHandler is the http request handler
// for posting of Game matchmaking
func gameHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	con := s.pool.Get()
	defer con.Close()

	log.Print("[Info][game_route] Match to a game")
	g, err := popOpenGame(con)

	if err != nil {
		if err != errGameNotFound {
			return err
		}

		g = NewGame()
		err := pushOpenGame(con, g)

		if err != nil {
			return err
		}

		// return 201 when pushing into the list
		w.WriteHeader(http.StatusCreated)
	} else {
		// creates the running server, and returns
		// a game with the ip and port populated
		g, err = s.createSession(con, g)
		if err != nil {
			return err
		}

		// update the record
		err := updateGame(con, g)
		if err != nil {
			return err
		}
	}

	err = json.NewEncoder(w).Encode(g)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		return err
	}

	return nil
}
