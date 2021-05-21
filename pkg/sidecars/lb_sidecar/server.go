package sidecarserv

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/inspr/inspr/pkg/sidecar_old/models"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	writeAddr string
	readAddr  string
}

// NewServer returns a new sidecar server
func NewServer() *Server {
	return &Server{}
}

// Init - configures the server
func (s *Server) Init(r models.Reader, w models.Writer) {
	// server requests related
	s.writeAddr = fmt.Sprintf(":%s", os.Getenv("INSPR_LBSIDECAR_WRITE_PORT"))
	s.readAddr = fmt.Sprintf(":%s", os.Getenv("INSPR_LBSIDECAR_READ_PORT"))
}

// Run starts the server on the port given in addr
func (s *Server) Run(ctx context.Context) {
	errCh := make(chan error)

	writeServer := &http.Server{
		Handler: s.writeMessageHandler().Post().JSON(),
		Addr:    s.writeAddr,
	}
	go func(ctx context.Context) {
		if err := writeServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			logger.Sugar().Fatalf("listen:%v", err)
		}
	}(ctx)

	readServer := &http.Server{
		Handler: s.readMessageHandler().Post().JSON(),
		Addr:    s.readAddr,
	}
	go func(ctx context.Context) {
		if err := readServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			logger.Sugar().Fatalf("listen:%v", err)
		}
	}(ctx)

	logger.Info("LB Sidecar listener is up...")

	select {
	case <-ctx.Done():
		gracefulShutdown(writeServer, readServer, nil)
	case errRead := <-errCh:
		gracefulShutdown(writeServer, readServer, errRead)
	}

}

func gracefulShutdown(w, r *http.Server, err error) {
	logger.Info("gracefully shutting down...")

	ctxShutdown, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*5),
	)
	defer cancel()

	if err != nil {
		logger.Sugar().Fatalf(err.Error())
	}

	// has to be the last method called in the shutdown
	if err = w.Shutdown(ctxShutdown); err != nil {
		logger.Sugar().Fatalf("error shutting down server")
	}

	if err = r.Shutdown(ctxShutdown); err != nil {
		logger.Sugar().Fatalf("error shutting down server")
	}
}
