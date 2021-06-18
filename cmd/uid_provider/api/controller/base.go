package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"inspr.dev/inspr/cmd/uid_provider/client"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	mux *http.ServeMux
	rdb client.RedisManager
	ctx context.Context
}

// Init - configures the server
func (s *Server) Init(ctx context.Context, rdb client.RedisManager) {
	s.mux = http.NewServeMux()
	s.rdb = rdb
	s.ctx = ctx
	s.initRoutes()
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	fmt.Printf("uidp rest api is up! listening on port: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.mux))
}
