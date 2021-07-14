package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/uid_provider/client"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "uidp-api-controllers")))
}

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	mux *http.ServeMux
	rdb client.RedisManager
	ctx context.Context
}

// Init - configures the server
func (s *Server) Init(ctx context.Context, rdb client.RedisManager) {
	logger.Info("initializing UIDP server")

	s.mux = http.NewServeMux()
	s.rdb = rdb
	s.ctx = ctx
	s.initRoutes()
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	logger.Info("running UIDP server",
		zap.String("on address", addr))

	fmt.Printf("uidp rest api is up! listening on port: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.mux))
}
