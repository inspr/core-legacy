package dappclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

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

func (c *Client) WriteMessage(channel string, msg models.Message) error {
	reqBody := struct {
		channel string         `json:"channel"`
		msg     models.Message `json:"message"`
	}{channel, msg}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := c.httpc.Post(c.addr+"/writeMessage", "", bytes.NewBuffer(reqBytes))
	if err != nil {
		return errors.New("folder/routes doesn't exists")
	}

	if resp.StatusCode != http.StatusOK {
		return rest.UnmarshalERROR(resp.Body)
	}
	return nil
}

func (c *Client) ReadMessage(channel string) (models.Message, error) {
	reqBody := struct {
		channel string
	}{channel}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return models.Message{}, err
	}

	resp, err := c.httpc.Post(c.addr+"/readMessage", "", bytes.NewBuffer(reqBytes))
	if err != nil {
		return models.Message{}, errors.New("folder/routes doesn't exists")
	}

	if resp.StatusCode != http.StatusOK {
		return models.Message{}, rest.UnmarshalERROR(resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	msg := models.Message{}
	decoder.Decode(&msg)

	return msg, nil
}

func (c *Client) CommitMessage(channel string) error {
	reqBody := struct {
		channel string
	}{channel}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := c.httpc.Post(c.addr+"/commitMessage", "", bytes.NewBuffer(reqBytes))
	if err != nil {
		return errors.New("folder/route doesn't exist")
	}

	if resp.StatusCode != http.StatusOK {
		return rest.UnmarshalERROR(resp.Body)
	}

	return nil
}
