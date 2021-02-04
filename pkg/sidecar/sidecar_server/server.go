package sidecarserv

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Mux *http.ServeMux
	sync.Mutex
	Reader models.Reader
	Writer models.Writer
}

// Init - configures the server
func (s *Server) Init(r models.Reader, w models.Writer) {
	// server requests related
	s.Mux = http.NewServeMux()

	// implementations of write and read for a specific sidecar
	s.Reader = r
	s.Writer = w

	s.InitRoutes()
}

// InitRoutes establishes the routes of the server
func (s *Server) InitRoutes() {
	handler := newCustomHandlers(s)

	s.Mux.HandleFunc("/writeMessage", handler.writeMessageHandler)

	s.Mux.HandleFunc("/readMessage", handler.readMessageHandler)
	s.Mux.HandleFunc("/commit", handler.commitMessageHandler)
}

// Run starts the server on the port given in addr
func (s *Server) Run(ctx context.Context) {
	server := &http.Server{
		Handler: s.Mux,
	}

	// todo: check if it can listen to the unix socket, os -> folder exists
	// logrus.Infoln("running write")
	// if _, err := os.Stat(WriteAddress); !os.IsNotExist(err) {
	// 	os.RemoveAll(WriteAddress)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	listenerAddr := environment.GetEnvironment().UnixSocketAddr

	// todo: replace this with a interface that returns a listener
	listener, err := net.Listen("unix", listenerAddr)
	if err != nil {
		log.Println("couldn't listen to address: " + listenerAddr)
		return
	}

	go func() {
		if err = server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%v", err)
		}
	}()

	fmt.Printf("SideCar listener is up...")
	<-ctx.Done()
	log.Printf("Gracefully shutting down...")

	ctxShutdown, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
	defer func() {
		cancel()
	}()

	if err = server.Shutdown(ctxShutdown); err != nil {
		log.Fatal("error shutting down server")
	}

	err = os.RemoveAll(listenerAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server shutdown complete")
	if err == http.ErrServerClosed {
		err = nil
	}
	return
}
