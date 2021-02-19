package controller

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Mux             *http.ServeMux
	MemoryManager   memory.Manager
	ChannelOperator operators.ChannelOperatorInterface
	NodeOperator    operators.NodeOperatorInterface
}

// Init - configures the server
func (s *Server) Init(mm memory.Manager, nodeOperator operators.NodeOperatorInterface, channelOperator operators.ChannelOperatorInterface) {
	s.Mux = http.NewServeMux()
	s.MemoryManager = mm
	s.ChannelOperator = channelOperator
	s.NodeOperator = nodeOperator
	s.initRoutes()
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	fmt.Printf("insprd rest api is up! Listening on port: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.Mux))
}
