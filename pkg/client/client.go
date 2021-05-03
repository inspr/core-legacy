package dappclient

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/inspr/inspr/pkg/rest/request"
	"github.com/inspr/inspr/pkg/sidecar/models"
	"github.com/inspr/inspr/pkg/sidecar/transports"
)

// Client is the struct which implements the methods of AppClient interface
type Client struct {
	client *request.Client
}

// clientMessage is the struct that represents the client's request format
type clientMessage struct {
	Message models.Message `json:"message"`
	Channel string         `json:"channel"`
}

// NewAppClient returns a new instance of the client of the AppClient package
func NewAppClient() *Client {
	socket := os.Getenv("INSPR_UNIX_SOCKET")
	if socket == "" {
		panic("NO SOCKET ENVIRONMENT VARIABLE")
	}
	envAddr := "/inspr/" + socket + ".sock"
	return &Client{
		client: request.NewClient().
			BaseURL("http://unix").
			HTTPClient(transports.NewUnixSocketClient(envAddr)).
			Encoder(json.Marshal).
			Decoder(request.JSONDecoderGenerator),
	}
}

// WriteMessage receives a channel and a message and sends it in a request to the sidecar server
func (c *Client) WriteMessage(ctx context.Context, channel string, msg models.Message) error {
	data := clientMessage{
		Channel: channel,
		Message: msg,
	}

	var resp interface{}
	err := c.client.Send(ctx, "/writeMessage", http.MethodPost, data, &resp)
	return err
}

// ReadMessage receives a channel and sends it in a request to the sidecar server
func (c *Client) ReadMessage(
	ctx context.Context,
	channel string,
	message interface{},
) error {
	data := clientMessage{
		Channel: channel,
	}

	err := c.client.Send(
		ctx,
		"/readMessage",
		http.MethodPost,
		data,
		message,
	)

	return err
}

// CommitMessage receives a channel and sends it in a request to the sidecar server
func (c *Client) CommitMessage(ctx context.Context, channel string) error {
	data := clientMessage{
		Channel: channel,
	}

	var resp interface{}
	err := c.client.Send(ctx, "/commit", http.MethodPost, data, &resp)

	return err
}
