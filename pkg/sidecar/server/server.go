package sidecarserv

import (
	"context"
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
	handler := newCustomHandlers(&s.Mutex, s.Reader, s.Writer)

	s.Mux.HandleFunc("/writeMessage", handler.writeMessageHandler)

	s.Mux.HandleFunc("/readMessage", handler.readMessageHandler)
	s.Mux.HandleFunc("/commit", handler.commitMessageHandler)
}

// Run starts the server on the port given in addr
func (s *Server) Run(ctx context.Context) {
	server := &http.Server{
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           s.Mux,
	}

	listenerAddr := environment.GetEnvironment().UnixSocketAddr
	os.Remove(listenerAddr)

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

	log.Printf("sideCar listener is up...")
	select {
	case <-ctx.Done():
		log.Println("gracefully shutting down...")

		ctxShutdown, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
		defer cancel()

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
}
