package dappclient

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/transports"
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

// requestReturn is the struct that represents the sidecar server's response
type requestReturn struct {
	Error   error          `json:"error"`
	Message models.Message `json:"message"`
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
			Decoder(request.JSONDecoderGenerator).
			Build(),
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

	// if reflect.ValueOf(message).Kind() != reflect.Struct {
	// 	return ierrors.NewError().
	// 		Message("message was not a struct").
	// 		Build()
	// }

	err := c.client.Send(
		ctx,
		"/readMessage",
		http.MethodPost,
		data,
		&message,
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
