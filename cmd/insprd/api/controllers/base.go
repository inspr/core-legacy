package controller

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Mux           *http.ServeMux
	MemoryManager memory.Manager
}

// Init todo doc
func (s *Server) Init() {
	s.Mux = http.NewServeMux()
	s.initRoutes()
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	fmt.Printf("insprd rest api is up! Listening on port: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.Mux))
}
