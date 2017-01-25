package sessions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// RegisterHandler registration of a new game session
func RegisterHandler(s *Server, w http.ResponseWriter, r *http.Request) error {
	var sess Session

	//TODO: this needs cleaning up allot

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

	cs, err := ClientSet()

	if err != nil {
		log.Printf("[Error][Register] Error connecting to Kubernetes: %v", err)
		return err
	}

	list, err := HostNameAndIP(cs)

	if err != nil {
		log.Printf("[Error][Register] Error Listing nodes: %v", err)
		return err
	}

	ip, err := ExternalNodeIPofPod(cs, sess, list)
	if err != nil {
		return err
	}

	log.Printf("[Info][Register] Session: IP: %v, Port: %v", ip, sess.Port)

	return s.storeSession(sess)
}
