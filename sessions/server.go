package sessions

import (
	"log"
	"net/http"

	"time"

	"github.com/cenkalti/backoff"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

// Version is the current api version number
const Version string = "0.2"

// Server is the http server instance
type Server struct {
	addr string
	pool *redis.Pool
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

	s := &Server{addr: hostAddr, pool: newPool(redisAddr)}
	return s
}

// Start starts the HTTP server on the given port
func (s *Server) Start() {

	r := mux.NewRouter()
	r.HandleFunc("/register", s.standardHandler(RegisterHandler)).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    s.addr,
	}

	err := s.pingRedis()
	if err != nil {
		log.Fatalf("[Error][Server] Could not connect to redis: %v", err)
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

// pingRedis pings redis, to check if we are
// connected. Returns an error if there was a problem.
func (s *Server) pingRedis() error {

	return backoff.Retry(func() error {
		con := s.pool.Get()
		defer con.Close()

		_, err := con.Do("PING")
		if err != nil {
			log.Printf("[Warn][Redis] Could not connect to Redis. %v", err)
		} else {
			log.Print("[Info][Redis] Connected.")
		}

		return err
	}, backoff.NewExponentialBackOff())
}

func newPool(address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
