package sessions

import (
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/markmandel/paddle-soccer/pkg"
	"k8s.io/client-go/kubernetes"
)

// Version is the current api version number
const Version string = "0.3"

// Server is the http server instance
type Server struct {
	addr string
	pool *redis.Pool
	cs   kubernetes.Interface
}

// Handler is the extended http.HandleFunc to provide context for this application
type Handler func(*Server, http.ResponseWriter, *http.Request) error

// NewServer returns the HTTP Server instance
func NewServer(hostAddr, redisAddr string) *Server {
	if redisAddr == "" {
		redisAddr = ":6379"
	}

	log.Printf("[Info][Server] Starting server version %v on port %v", Version, hostAddr)
	log.Printf("[Info][Server] Connecting to Redis at %v", redisAddr)

	s := &Server{addr: hostAddr, pool: pkg.NewPool(redisAddr)}
	return s
}

// Start starts the HTTP server on the given port
func (s *Server) Start() {

	r := mux.NewRouter()
	r.HandleFunc("/register", s.standardHandler(registerHandler)).Methods("POST")
	r.HandleFunc("/session/{id}", s.standardHandler(getHandler)).Methods("GET")
	r.HandleFunc("/session", s.standardHandler(createHandler)).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    s.addr,
	}

	err := pkg.WaitForConnection(s.pool)
	if err != nil {
		log.Fatalf("[Error][Server] Could not connect to redis: %v", err)
	}

	s.cs, err = clientSet()
	if err != nil {
		log.Fatalf("[Error][Server] Could not connect to kubernetes api: %v", err)
	}

	log.Fatalf("[Error][Server] Error starting server: %v", srv.ListenAndServe())
}

// standardHandler returns a http.HandleFunc
// wrapped in standard middleware
func (s *Server) standardHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(s, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
