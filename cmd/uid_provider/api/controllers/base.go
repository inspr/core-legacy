package controller

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/uid_provider/client"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	mux *http.ServeMux
	rdb client.RedisManager
}

// Init - configures the server
func (s *Server) Init(rdb client.RedisManager) {
	s.mux = http.NewServeMux()
	s.rdb = rdb
	s.initRoutes()
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	fmt.Printf("insprd rest api is up! Listening on port: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.mux))
}
