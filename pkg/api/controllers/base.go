package controller

import (
	"log"
	"net/http"

	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/operators"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/logs"
)

var logger *zap.Logger
var alevel *zap.AtomicLevel

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, alevel = logs.Logger(
		zap.Fields(zap.String("section", "insprd-api-controllers")),
	)
}

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	mux    *http.ServeMux
	memory memory.Manager
	op     operators.OperatorInterface
	auth   auth.Auth
}

// Init - configures the server
func (s *Server) Init(
	mem memory.Manager,
	op operators.OperatorInterface,
	auth auth.Auth,
) {
	logger.Info("initializing Insprd server")

	s.mux = http.NewServeMux()
	s.memory = mem
	s.op = op
	s.auth = auth
	s.initRoutes()
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	logger = logger.With(zap.String("port", addr))
	logger.Info("running Insprd server")

	log.Fatal(http.ListenAndServe(addr, s.mux))
}
