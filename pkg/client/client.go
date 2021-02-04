package dappclient

// todo check if Post methods are capable of receiving contexts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/transports"
)

// Client struct that implements the methods of the AppClient interface
type Client struct {
	addr  string
	httpc http.Client
}

// NewAppClient returns a new instance of the client of the AppClient package
func NewAppClient() *Client {
	// todo get env var
	envAddr := "/"
	return &Client{
		addr:  envAddr,
		httpc: transports.NewUnixSocketClient(envAddr),
	}
}

// WriteMessage TODO DOC
func (c *Client) WriteMessage(channel string, msg models.Message) error {
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)

	go func() {
		defer cancel()
		// todo: define struct elsewhere
		reqBody := struct {
			channel string         `json:"channel"`
			msg     models.Message `json:"message"`
		}{channel, msg}

		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			errChan <- err
			return
		}

		resp, err := c.httpc.Post(c.addr+"/writeMessage", "", bytes.NewBuffer(reqBytes))
		if err != nil {
			errChan <- errors.New("folder/routes doesn't exists")
			return
		}

		if resp.StatusCode != http.StatusOK {
			errChan <- rest.UnmarshalERROR(resp.Body)
			return
		}
		errChan <- nil
	}()

	for {
		select {
		case wmErr := <-errChan:
			return wmErr

		case <-ctx.Done():
			return ierrors.NewError().InternalServer().Message("server died mid write message request").Build()
		}
	}
}

// ReadMessage TODO DOC
func (c *Client) ReadMessage(channel string) (models.Message, error) {
	// todo: define struct elsewhere
	type rmReturns struct {
		Error   error
		Message models.Message
	}

	ctx, cancel := context.WithCancel(context.Background())
	ret := make(chan rmReturns)

	go func() {
		defer cancel()
		reqBody := struct {
			channel string
		}{channel}

		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			ret <- rmReturns{err, models.Message{}}
			return
		}

		resp, err := c.httpc.Post(c.addr+"/readMessage", "", bytes.NewBuffer(reqBytes))
		if err != nil {
			ret <- rmReturns{errors.New("folder/routes doesn't exists"), models.Message{}}
			return
		}

		if resp.StatusCode != http.StatusOK {
			ret <- rmReturns{rest.UnmarshalERROR(resp.Body), models.Message{}}
			return
		}

		decoder := json.NewDecoder(resp.Body)
		msg := models.Message{}
		decoder.Decode(&msg)
		ret <- rmReturns{nil, msg}
	}()

	for {
		select {
		case rmErr := <-ret:
			return rmErr.Message, rmErr.Error
		case <-ctx.Done():
			return models.Message{}, ierrors.NewError().InternalServer().Message("server died mid read message request").Build()
		}
	}
}

// CommitMessage TODO DOC
func (c *Client) CommitMessage(channel string) error {
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)

	go func() {
		defer cancel()
		reqBody := struct {
			channel string
		}{channel}

		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			errChan <- err
			return
		}

		resp, err := c.httpc.Post(c.addr+"/commitMessage", "", bytes.NewBuffer(reqBytes))
		if err != nil {
			errChan <- errors.New("folder/route doesn't exist")
			return
		}

		if resp.StatusCode != http.StatusOK {
			errChan <- rest.UnmarshalERROR(resp.Body)
			return
		}
		errChan <- nil
	}()

	for {
		select {
		case wmErr := <-errChan:
			return wmErr
		case <-ctx.Done():
			return ierrors.NewError().InternalServer().Message("server died mid commit message request").Build()
		}
	}
}
