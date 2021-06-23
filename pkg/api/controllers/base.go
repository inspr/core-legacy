package controller

import (
	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
	"inspr.dev/inspr/cmd/insprd/operators"
	"inspr.dev/inspr/pkg/auth"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "insprd-api-controllers")))
}

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Mux               *http.ServeMux
	TreeMemoryManager tree.Manager
	BrokerManager     brokers.Manager
	op                operators.OperatorInterface
	auth              auth.Auth
}

// Init - configures the server
func (s *Server) Init(mm tree.Manager, op operators.OperatorInterface, auth auth.Auth, bm brokers.Manager) {
	logger.Info("initializing Insprd server")

	s.Mux = http.NewServeMux()
	s.TreeMemoryManager = mm
	s.BrokerManager = bm
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
