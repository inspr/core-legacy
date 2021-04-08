package controllers

import (
	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Mux    *http.ServeMux
	logger *zap.Logger
}

// Init - configures the server
func (s *Server) Init() {
	s.Mux = http.NewServeMux()
	s.logger, _ = zap.NewDevelopment(zap.Fields(zap.String("section", "Auth-provider")))
	s.initRoutes()
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	fmt.Printf("insprd rest api is up! Listening on port: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.Mux))
}
