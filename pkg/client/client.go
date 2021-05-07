package dappclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/rest/request"
)

// Client is the struct which implements the methods of AppClient interface
type Client struct {
	client   *request.Client
	mux      *http.ServeMux
	readAddr string
}

// clientMessage is the struct that represents the client's request format
type clientMessage struct {
	Message interface{} `json:"message"`
}

// NewAppClient returns a new instance of the client of the AppClient package
func NewAppClient() *Client {

	writeAddr := fmt.Sprintf("http://localhost:%s", os.Getenv("INSPR_SIDECAR_WRITE_PORT"))
	readAddr := fmt.Sprintf(":%s", os.Getenv("INSPR_SIDECAR_READ_PORT"))
	return &Client{
		readAddr: readAddr,
		client: request.NewClient().
			BaseURL(writeAddr).
			Encoder(json.Marshal).
			Decoder(request.JSONDecoderGenerator).
			Build(),
		mux: http.NewServeMux(),
	}
}

// WriteMessage receives a channel and a message and sends it in a request to the sidecar server
func (c *Client) WriteMessage(ctx context.Context, channel string, msg interface{}) error {
	data := clientMessage{
		Message: msg,
	}

	var resp interface{}
	log.Println("sending message to sidecar")
	// sends a message to the corresponding channel route on the sidecar
	err := c.client.Send(ctx, "/"+channel, http.MethodPost, data, &resp)
	log.Println("message sent")
	return err
}

// HandleChannel handles messages received in a given channel.
func (c *Client) HandleChannel(channel string, handler func(ctx context.Context, body io.Reader) error) {
	c.mux.HandleFunc("/"+channel, func(w http.ResponseWriter, r *http.Request) {
		// user defined handler. Returns error if the user wants to return it
		err := handler(context.Background(), r.Body)
		if err != nil {
			rest.ERROR(w, ierrors.NewError().InternalServer().InnerError(err).Build())
			return
		}
		rest.JSON(w, 200, struct{ Status string }{"OK"})

	})
}

//Run runs the server with the handlers defined in HandleChannel
func (c *Client) Run(ctx context.Context) error {

	var err error
	server := http.Server{
		Handler: c.mux,
		Addr:    c.readAddr,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%v", err)
		}
	}()

	log.Printf("sideCar listener is up...")

	<-ctx.Done()

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
		return err
	}
	return ctx.Err()
}
