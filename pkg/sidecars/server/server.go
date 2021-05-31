package sidecarserv

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/inspr/inspr/pkg/sidecars/models"
	"go.uber.org/zap"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	broker       string
	Reader       models.Reader
	Writer       models.Writer
	inAddr       string
	outAddr      string
	client       *http.Client
	runningRead  bool
	runningWrite bool
}

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "sidecar")))
}

// Init - configures the server
func Init(r models.Reader, w models.Writer, broker string) *Server {
	server := &Server{}
	// server fetches addresses variable names from models.
	envVars := models.GetSidecarConnectionVars(broker)
	if envVars == nil {
		panic(fmt.Sprintf("%s broker's enviroment variables not configured", broker))
	}
	server.broker = broker

	// server fetches required addresses from deployment.
	inAddr, ok := os.LookupEnv(envVars.WriteEnvVar)
	if !ok {
		panic(fmt.Sprintf("[ENV VAR] %s not found", envVars.WriteEnvVar))
	}

	outAddr, ok := os.LookupEnv(envVars.ReadEnvVar)
	if !ok {
		panic(fmt.Sprintf("[ENV VAR] %s not found", envVars.ReadEnvVar))
	}

	server.inAddr = fmt.Sprintf(":%s", inAddr)
	server.outAddr = fmt.Sprintf("http://localhost:%v", outAddr)
	server.client = &http.Client{}

	// implementations of write and read for a specific sidecar
	server.Reader = r
	server.Writer = w
	return server
}

// Run starts the server on the port given in addr
func (s *Server) Run(ctx context.Context) error {
	server := &http.Server{
		Handler: s.writeMessageHandler().Post().JSON(),
		Addr:    s.inAddr,
	}
	errCh := make(chan error)
	// create read message routine and captures its error
	go func() { errCh <- s.readMessageRoutine(ctx) }()

	var err error
	go func() {
		s.runningWrite = true
		defer func() { s.runningWrite = false }()
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen:%v", zap.Error(err))
			errCh <- err
		}
	}()

	logger.Sugar().Infof("%s sideCar listener is up...", s.broker)

	select {
	case <-ctx.Done():
		s.gracefulShutdown(server, ctx.Err())
		return ctx.Err()
	case errRead := <-errCh:
		s.gracefulShutdown(server, errRead)
		return errRead
	}

}

func (s *Server) gracefulShutdown(server *http.Server, err error) {
	logger.Info("gracefully shutting down...")

	ctxShutdown, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*5),
	)

	defer cancel()

	if err != nil {
		logger.Error("an error occurred on sidecar", zap.Any("broker", s.broker), zap.Error(err))
	}

	s.Writer.Close()

	// has to be the last method called in the shutdown
	if err = server.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("error shutting down server")
	}
}
