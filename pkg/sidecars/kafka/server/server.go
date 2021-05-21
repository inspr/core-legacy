package sidecarserv

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/inspr/inspr/pkg/rest/request"
	"github.com/inspr/inspr/pkg/sidecar_old/models"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Reader       models.Reader
	Writer       models.Writer
	writeAddr    string
	client       *request.Client
	runningRead  bool
	runningWrite bool
}

// NewServer returns a new sidecar server
func NewServer() *Server {
	return &Server{}
}

// Init - configures the server
func (s *Server) Init(r models.Reader, w models.Writer) {

	// server fetches required addresses from deployment.

	wAddr, ok := os.LookupEnv("INSPR_SIDECAR_KAFKA_WRITE_PORT")
	if !ok {
		panic("[ENV VAR] INSPR_SIDECAR_KAFKA_WRITE_PORT not found")
	}

	rAddr, ok := os.LookupEnv("INSPR_SIDECAR_KAFKA_READ_PORT")
	if !ok {
		panic("[ENV VAR] INSPR_SIDECAR_KAFKA_READ_PORT not found")
	}

	s.writeAddr = fmt.Sprintf(":%s", wAddr)
	s.client = request.NewJSONClient(fmt.Sprintf("http://localhost:%v", rAddr))

	// implementations of write and read for a specific sidecar
	s.Reader = r
	s.Writer = w
} // looked

// Run starts the server on the port given in addr
func (s *Server) Run(ctx context.Context) {
	server := &http.Server{
		Handler: s.writeMessageHandler().Post().JSON(), // look writeMessageHandler
		Addr:    s.writeAddr,
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
		gracefulShutdown(server, err)
	case errRead := <-errCh:
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
