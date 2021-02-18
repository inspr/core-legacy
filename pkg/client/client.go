package dappclient

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/environment"
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
	envAddr := "/inspr/" + environment.GetEnvironment().UnixSocketAddr + ".sock"
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
func (client *Client) WriteMessage(ctx context.Context, channel string, msg models.Message) error {
	data := clientMessage{
		Channel: channel,
		Message: msg,
	}

	var resp interface{}

	err := client.client.Send(ctx, "/writeMessage", http.MethodPost, data, &resp)
	return err
}

// ReadMessage receives a channel and sends it in a request to the sidecar server
func (client *Client) ReadMessage(ctx context.Context, channel string) (models.Message, error) {
	data := clientMessage{
		Channel: channel,
	}

	var msg models.BrokerData

	err := client.client.Send(ctx, "/readMessage", http.MethodPost, data, &msg)
	return msg.Message, err
}

// CommitMessage receives a channel and sends it in a request to the sidecar server
func (client *Client) CommitMessage(ctx context.Context, channel string) error {
	data := clientMessage{
		Channel: channel,
	}

	var resp interface{}

	err := client.client.Send(ctx, "/commit", http.MethodPost, data, &resp)

	return err
}
