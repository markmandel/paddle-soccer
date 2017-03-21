package nodescaler

import (
	"log"
)

// scale scales nodes up and down, depending on CPU constraints
func (s Server) scaleNodes() error {
	log.Print("[Debug][Scaler] ...Checking Scale... ")
	nl, err := s.newNodeList()
	if err != nil {
		return err
	}

	available := s.cpuRequestsAvailable(nl)
	log.Printf("[Debug][Scaler] CPU Requests blocks of %v available: %v", s.cpuRequest, available)

	return nil
}
