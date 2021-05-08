package client

import (
	"encoding/json"
	"os"

	"github.com/inspr/inspr/pkg/controller"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest/request"
)

const inClusterEnviromentError = "authentication as controller failed. controllers requires following " +
	"variables: INSPR_INSPRD_ADDRESS, INSPR_CONTROLLER_SCOPE & INSPR_CONTROLLER_TOKEN"

// ControllerConfig stores controller configuration for ease of use and posterior verification.
type ControllerConfig struct {
	Auth  request.Authenticator
	Scope string
	URL   string
}

type authenticator struct{}

func (authenticator) GetToken() ([]byte, error) {
	return []byte("Bearer " + os.Getenv("INSPR_CONTROLLER_TOKEN")), nil
}
func (authenticator) SetToken(tok []byte) error {
	os.Setenv("INSPR_CONTROLLER_TOKEN", string(tok)[len("Bearer "):])
	return nil
}

// GetInClusterConfigs retrieves controller configs from current dApp deployment.
func GetInClusterConfigs() (*ControllerConfig, error) {
	url, urlok := os.LookupEnv("INSPR_INSPRD_ADDRESS")
	scope, scopeok := os.LookupEnv("INSPR_CONTROLLER_SCOPE")
	_, tknok := os.LookupEnv("INSPR_CONTROLLER_TOKEN")
	if !urlok || !scopeok || !tknok {
		return nil, ierrors.NewError().
			Message(inClusterEnviromentError).
			Build()
	}
	return &ControllerConfig{
		Auth:  authenticator{},
		Scope: scope,
		URL:   url,
	}, nil
}

// Client implements communication with the Insprd
type Client struct {
	HTTPClient *request.Client
	Config     ControllerConfig
}

// NewControllerClient return a new Client
func NewControllerClient(config ControllerConfig) controller.Interface {
	return &Client{
		HTTPClient: request.NewClient().
			BaseURL(config.URL).
			Encoder(json.Marshal).
			Decoder(request.JSONDecoderGenerator).
			Authenticator(config.Auth).
			Pointer(),
	}
}

// Channels interacts with channels on the Insprd
func (c *Client) Channels() controller.ChannelInterface {
	return &ChannelClient{
		reqClient: c.HTTPClient,
	}
}

// Apps interacts with apps on the Insprd
func (c *Client) Apps() controller.AppInterface {
	return &AppClient{
		reqClient: c.HTTPClient,
	}
}

// ChannelTypes interacts with channel types on the Insprd
func (c *Client) ChannelTypes() controller.ChannelTypeInterface {
	return &ChannelTypeClient{
		reqClient: c.HTTPClient,
	}
}

// Authorization interacts with Insprd's auth
func (c *Client) Authorization() controller.AuthorizationInterface {
	return &AuthClient{
		reqClient: c.HTTPClient,
	}
}

// Alias interacts with alias on the Insprd
func (c *Client) Alias() controller.AliasInterface {
	return &AliasClient{
		reqClient: c.HTTPClient,
	}
}
