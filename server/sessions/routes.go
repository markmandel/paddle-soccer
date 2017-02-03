package sessions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// createHandler registration of a new game session
func createHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	log.Printf("[Info][create] creating a game session with image: %v", s.gameServerImage)

	id, err := s.createSessionPod()
	if err != nil {
		return err
	}

	sess := Session{ID: id}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(&sess)
}

// RegisterHandler registration of a new game session
func registerHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	var sess Session

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[Error][Register] Reading Body: %v.", err)
		return err
	}

	log.Printf("[Info][Register] Recieved JSON Playload: %v", string(b))

	err = json.Unmarshal(b, &sess)

	if err != nil {
		log.Printf("[Error][Register] Error decoding json: %v. [%v]", err, string(b))
		return err
	}

	log.Printf("[Info][Register] Recieved Session Registration: %#v", sess)

	if err != nil {
		log.Printf("[Error][Register] Error connecting to Kubernetes: %v", err)
		return err
	}

	list, err := s.hostNameAndIP()
	if err != nil {
		return err
	}

	ip, err := s.externalNodeIPofPod(sess, list)
	if err != nil {
		return err
	}

	log.Printf("[Info][Register] Session: IP: %v, Port: %v", ip, sess.Port)
	sess.IP = ip

	return s.storeSession(sess)
}

// getHandler for retrieving information about a game session
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
