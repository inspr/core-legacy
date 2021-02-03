package sidecarserv

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"gitlab.inspr.dev/inspr/core/pkg/sidecar"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Mux         *http.ServeMux
	Mutex       *sync.Mutex
	Messages    chan sidecar.Message
	receiveAddr string
	sendAddr    string
}

// Init - configures the server
func (s *Server) Init(receiveAddr, sendAddr string) {

	// server requests related
	s.Mux = http.NewServeMux()
	s.Mutex = &sync.Mutex{}

	// listener and destination routes
	s.receiveAddr = receiveAddr
	s.sendAddr = sendAddr

	// limits the amount of messages to 100
	s.Messages = make(chan sidecar.Message, 100)

	s.InitRoutes()
}

// InitRoutes establishes the routes of the server
func (s *Server) InitRoutes() {
	customHandlers := newHandlers(s)
	s.Mux.HandleFunc("/message", customHandlers.messageHandler)
	s.Mux.HandleFunc("/commmit", customHandlers.commitHandler)

}

// Run starts the server on the port given in addr
func (s *Server) Run() {
	fmt.Printf("SideCart listerner is up...")
	// todo: replace this with a interface that returns a listener
	listen, _ := net.Listen("unix", s.receiveAddr)
	http.Serve(listen, s.Mux)
}
