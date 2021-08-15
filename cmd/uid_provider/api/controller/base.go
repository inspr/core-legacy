package controller

import (
	"context"
	"log"
	"net/http"

	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/uid_provider/client"
	"inspr.dev/inspr/pkg/logs"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = logs.Logger(
		zap.Fields(zap.String("section", "uidp-api-controllers")),
	)
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
	logger = logger.With(zap.String("port", addr))
	logger.Info("listening")

	log.Fatal(http.ListenAndServe(addr, s.mux))
}
