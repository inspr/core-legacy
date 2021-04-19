package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/operators"
	"github.com/inspr/inspr/pkg/auth"
	"go.uber.org/zap"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = zap.NewDevelopment(zap.Fields(zap.String("section", "insprd-api-controllers")))
}

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Mux           *http.ServeMux
	MemoryManager memory.Manager
	op            operators.OperatorInterface
	auth          auth.Auth
}

// Init - configures the server
func (s *Server) Init(mm memory.Manager, op operators.OperatorInterface, auth auth.Auth) {
	logger.Info("initializing Insprd server")

	s.Mux = http.NewServeMux()
	s.MemoryManager = mm
	s.op = op
	s.auth = auth
	s.initRoutes()
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	logger.Info("running Insprd server",
		zap.String("on address", addr))

	fmt.Printf("insprd rest api is up! Listening on port: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.Mux))
}
