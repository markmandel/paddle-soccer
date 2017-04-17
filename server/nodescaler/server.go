// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nodescaler

import (
	"log"
	"net/http"

	"time"

	"github.com/gorilla/mux"
	"github.com/markmandel/paddle-soccer/server/nodescaler/gce"
	"github.com/markmandel/paddle-soccer/server/pkg/kube"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/resource"
)

// Version is the current api version number
const Version string = "nodescaler:0.1"

// Server is the http server instance
type Server struct {
	srv *http.Server
	// `nodeSelector` is a k8s selector for what nodes to manage
	cs           kubernetes.Interface
	nodeSelector string
	// `cpuRequest` is the cpu capacity requested for each server
	cpuRequest int64
	// `bufferCount``is the number of cpuRequest to make sure is available
	// at and given moment in the nodePool
	bufferCount int64
	// nodePool management implementation.
	// for now, there is just GCE
	nodePool NodePool
	// minimum pool size
}

// handler is the extended http.HandleFunc to provide context for this application
type handler func(*Server, http.ResponseWriter, *http.Request) error

// NewServer returns the HTTP Server instance
// `nodeSelector` is a k8s selector for what nodes to manage
// `cpuRequest` is the cpu capacity requested for each server
// `bufferCount``is the number of cpuRequest to make sure is available
// at and given moment in the cluster
func NewServer(hostAddr, nodeSelector string, cpuRequest string, bufferCount int64) (*Server, error) {
	log.Printf("[Info][Server] Starting server version %v on port %v, managing node selector %v",
		Version, hostAddr, nodeSelector)

	q, err := resource.ParseQuantity(cpuRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not parse cpu resource request: %v", cpuRequest)
	}

	s := &Server{nodeSelector: nodeSelector, cpuRequest: q.MilliValue(), bufferCount: bufferCount}
	r := s.newHandler()

	s.srv = &http.Server{
		Handler: r,
		Addr:    hostAddr,
	}

	return s, nil
}

// Start starts the HTTP server on the given port
func (s *Server) Start() error {
	quit := make(chan bool)
	defer close(quit)

	var err error
	s.cs, err = kube.ClientSet()
	if err != nil {
		return errors.Wrap(err, "Could not connect to kubernetes api")
	}

	nl, err := s.newNodeList()
	if err != nil {
		return errors.WithMessage(err, "Could not create nodelist when starting Server")
	}
	// Hardcode to GCE for this proof of concept. Long term, this should be switchable.
	np, err := gce.NewNodePool(nl.nodes.Items[0])
	if err != nil {
		return err
	}
	s.nodePool = np

	// watch for the nodepool
	gw, err := s.newGameWatcher()
	if err != nil {
		return err
	}
	gw.start()

	go func() {
		log.Print("[Info][Start] Starting node scaling...")
		tick := time.Tick(10 * time.Second) //TODO: Make this an env var

		for {
			select {
			case <-quit:
				return
			case <-gw.event:
				log.Print("[Info][Scaling] Recieved Add Event, Scaling...")
				if err := s.scaleNodes(); err != nil {
					log.Printf("[Error][Scaling] %+v", err)
				}
			case <-tick:
				log.Printf("[Info][Scaling] Tick of %#v, Scaling...", tick)
				if err := s.scaleNodes(); err != nil {
					log.Printf("[Error][Scaling] %+v", err)
				}
			}
		}
	}()

	return errors.Wrap(s.srv.ListenAndServe(), "Error starting server")
}

// newHandler creates the http routes for this application
func (s *Server) newHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", func(http.ResponseWriter, *http.Request) {})

	return r
}

// wrapMiddleware returns a http.HandleFunc // wrapped in standard middleware
func (s *Server) wrapMiddleware(h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(s, w, r)
		if err != nil {
			log.Printf("[Error][Server] %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
