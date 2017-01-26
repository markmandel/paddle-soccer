package sessions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// RegisterHandler registration of a new game session
func registerHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	var sess Session

	//TODO: this needs cleaning up a lot

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
		log.Printf("[Error][Register] Error Listing nodes: %v", err)
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
