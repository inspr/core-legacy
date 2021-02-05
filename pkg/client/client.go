package dappclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/transports"
)

// Client is the struct which implements the methods of AppClient interface
type Client struct {
	addr  string
	httpc http.Client
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
	envAddr := environment.GetEnvironment().UnixSocketAddr
	return &Client{
		addr:  envAddr,
		httpc: transports.NewUnixSocketClient(envAddr),
	}
}

// WriteMessage receives a channel and a message and sends it in a request to the sidecar server
func (c *Client) WriteMessage(ctx context.Context, channel string, msg models.Message) error {
	data := clientMessage{
		Channel: channel,
		Message: msg,
	}
	_, err := c.sendRequest(ctx, http.MethodPost, c.addr+"/writeMessage", data)
	return err
}

// ReadMessage receives a channel and sends it in a request to the sidecar server
func (c *Client) ReadMessage(ctx context.Context, channel string) (models.Message, error) {
	data := clientMessage{
		Channel: channel,
	}
	msg, err := c.sendRequest(ctx, http.MethodPost, c.addr+"/readMessage", data)
	return msg, err
}

// CommitMessage receives a channel and sends it in a request to the sidecar server
func (c *Client) CommitMessage(ctx context.Context, channel string) error {
	data := clientMessage{
		Channel: channel,
	}
	_, err := c.sendRequest(ctx, http.MethodPost, c.addr+"/commit", data)
	return err
}

func (c *Client) sendRequest(ctx context.Context, method, addr string, reqData clientMessage) (models.Message, error) {
	ret := make(chan requestReturn)

	go func() {
		reqBytes, err := json.Marshal(reqData)
		if err != nil {
			ret <- requestReturn{err, models.Message{}}
			return
		}

		newRequest, err := http.NewRequest(method, addr, bytes.NewBuffer(reqBytes))
		if err != nil {
			ret <- requestReturn{err, models.Message{}}
			return
		}

		newRequest.WithContext(ctx)
		resp, err := c.httpc.Do(newRequest)
		if err != nil {
			ret <- requestReturn{err, models.Message{}}
			return
		}

		if resp.StatusCode != http.StatusOK {
			ret <- requestReturn{rest.UnmarshalERROR(resp.Body), models.Message{}}
			return
		}

		decoder := json.NewDecoder(resp.Body)
		msg := models.Message{}
		decoder.Decode(&msg)
		ret <- requestReturn{nil, msg}
	}()

	for {
		select {
		case rmErr := <-ret:
			return rmErr.Message, rmErr.Error
		case <-ctx.Done():
			return models.Message{}, ierrors.NewError().InternalServer().Message("server died mid request").Build()
		}
	}
}
