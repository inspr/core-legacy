package sidecarserv

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Mux *http.ServeMux
	sync.Mutex
	Messages   chan models.Message
	listenAddr string
	Reader     models.Reader
	Writer     models.Writer
}

// Init - configures the server
func (s *Server) Init(listenAddr string, r models.Reader, w models.Writer) {
	// server requests related
	s.Mux = http.NewServeMux()

	// listener and destination routes
	s.listenAddr = listenAddr

	// limits the amount of messages to 100
	s.Messages = make(chan models.Message, 100)

	// implementations of write and read for a specific sidecar
	s.Reader = r
	s.Writer = w

	s.InitRoutes()
}

// InitRoutes establishes the routes of the server
func (s *Server) InitRoutes() {
	s.Mux.HandleFunc("/writeMessage", s.writeMessageHandler)

	s.Mux.HandleFunc("/readMessage", s.readMessageHandler)
	s.Mux.HandleFunc("/commmit", s.commitMessageHandler)

}

// Run starts the server on the port given in addr
func (s *Server) Run() {
	// todo: check if it can listen to the unix socket, os -> folder exists
	// logrus.Infoln("running write")
	// if _, err := os.Stat(WriteAddress); !os.IsNotExist(err) {
	// 	os.RemoveAll(WriteAddress)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	fmt.Printf("SideCart listener is up...")
	// todo: replace this with a interface that returns a listener
	listen, _ := net.Listen("unix", s.listenAddr)
	http.Serve(listen, s.Mux)

	// todo: log of chan error
}
