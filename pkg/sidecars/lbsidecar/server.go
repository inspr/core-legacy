package lbsidecar

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	writeAddr string
	readAddr  string
}

// Init - initializes a new configured server
func Init() *Server {
	s := Server{}

	wAddr, exists := os.LookupEnv("INSPR_LBSIDECAR_WRITE_PORT")
	if !exists {
		panic("[ENV VAR] INSPR_LBSIDECAR_WRITE_PORT not found")
	}
	rAddr, exists := os.LookupEnv("INSPR_LBSIDECAR_READ_PORT")
	if !exists {
		panic("[ENV VAR] INSPR_LBSIDECAR_READ_PORT not found")
	}

	s.writeAddr = fmt.Sprintf(":%s", wAddr)
	s.readAddr = fmt.Sprintf(":%s", rAddr)
	return &s
}

// Run starts the server on the port given in addr
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error)

	writeServer := &http.Server{
		Handler: s.writeMessageHandler().Post().JSON(),
		Addr:    s.writeAddr,
	}
	go func() {
		if err := writeServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			logger.Error("an error ocurred in LB Sidecar write server",
				zap.Error(err))
		}
	}()

	readServer := &http.Server{
		Handler: s.readMessageHandler().Post().JSON(),
		Addr:    s.readAddr,
	}
	go func() {
		if err := readServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			logger.Error("an error ocurred in LB Sidecar read server",
				zap.Error(err))
		}
	}()

	logger.Info("LB Sidecar listener is up...")

	select {
	case <-ctx.Done():
		gracefulShutdown(writeServer, readServer, nil)
		return ctx.Err()
	case errRead := <-errCh:
		gracefulShutdown(writeServer, readServer, errRead)
		return errRead
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
		logger.Error("an error occured in LB Sidecar",
			zap.Error(err))
	}

	// has to be the last method called in the shutdown
	if err = w.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("error while shutting down LB Sidecar write server",
			zap.Error(err))
	}

	if err = r.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("error while shutting down LB Sidecar read server",
			zap.Error(err))
	}
}
