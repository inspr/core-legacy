package sidecarserv

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/inspr/inspr/pkg/sidecars/models"
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

// NewServer returns a new sidecar server
func NewServer() *Server {
	return &Server{}
}

// Init - configures the server
func (s *Server) Init(r models.Reader, w models.Writer, broker string) {

	// server fetches addresses variable names from models.
	envVars := models.GetSidecarConnectionVars(broker)
	if envVars == nil {
		panic(fmt.Sprintf("%s broker's enviroment variables not configured", broker))
	}
	s.broker = broker

	// server fetches required addresses from deployment.
	inAddr, ok := os.LookupEnv(envVars.WriteEnvVar)
	if !ok {
		panic(fmt.Sprintf("[ENV VAR] %s not found", envVars.WriteEnvVar))
	}

	outAddr, ok := os.LookupEnv(envVars.ReadEnvVar)
	if !ok {
		panic(fmt.Sprintf("[ENV VAR] %s not found", envVars.ReadEnvVar))
	}

	s.inAddr = fmt.Sprintf(":%s", inAddr)
	s.outAddr = fmt.Sprintf("http://localhost:%v", outAddr)
	s.client = &http.Client{}

	// implementations of write and read for a specific sidecar
	s.Reader = r
	s.Writer = w
} // looked

// Run starts the server on the port given in addr
func (s *Server) Run(ctx context.Context) {
	server := &http.Server{
		Handler: s.writeMessageHandler().Post().JSON(), // look writeMessageHandler - OK
		Addr:    s.inAddr,
	}
	errCh := make(chan error)
	// create read message routine and captures its error
	go func() { errCh <- s.readMessageRoutine(ctx) }() // readMessageRoutine

	var err error
	go func() {
		s.runningWrite = true
		defer func() { s.runningWrite = false }()
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%v", err)
		}
	}()

	log.Printf("sideCar listener is up...")

	select {
	case <-ctx.Done():
		s.Writer.Close()
		gracefulShutdown(server, err)
	case errRead := <-errCh:
		s.Writer.Close()
		gracefulShutdown(server, err)
		if errRead != nil {
			log.Fatalln(err)
		}
	}

}

func gracefulShutdown(server *http.Server, err error) {
	log.Println("gracefully shutting down...")

	ctxShutdown, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*5),
	)

	defer cancel()

	if err != nil {
		log.Fatal(err)
	}

	// has to be the last method called in the shutdown
	if err = server.Shutdown(ctxShutdown); err != nil {
		log.Fatal("error shutting down server")
	}
}
