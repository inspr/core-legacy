package client

import (
	"encoding/json"
	"os"

	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/controller"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest/request"
)

// Client implements communication with the Insprd
type Client struct {
	HTTPClient *request.Client
	Config     ControllerConfig
}

// ControllerConfig stores controller configuration for ease of use and posterior verification.
type ControllerConfig struct {
	Auth  request.Authenticator
	Scope string
	URL   string
}

// NewControllerClient return a new Client
func NewControllerClient(config ControllerConfig) controller.Interface {
	client := request.NewClient().BaseURL(config.URL).Encoder(json.Marshal).
		Decoder(request.JSONDecoderGenerator).Authenticator(config.Auth).Build()
	return &Client{
		HTTPClient: client,
		Config:     config,
	}
}

// GetInClusterConfigs retrieves controller configs from current dApp deployment.
func GetInClusterConfigs() (*ControllerConfig, error) {
	url, urlok := os.LookupEnv("INSPR_INSPRD_ADDRESS")
	scope, scopeok := os.LookupEnv("INSPR_CONTROLLER_SCOPE")
	token, tknok := os.LookupEnv("INSPR_CONTROLLER_TOKEN")
	if !urlok || !scopeok || !tknok {
		return nil, ierrors.NewError().
			Message("authentication as controller failed. controllers requires following variables: INSPR_INSPRD_ADDRESS, INSPR_CONTROLLER_SCOPE & INSPR_CONTROLLER_TOKEN").
			Build()
	}
	return &ControllerConfig{
		Auth: auth.Authenticator{
			TokenPath: token,
		},
		Scope: scope,
		URL:   url,
	}, nil
}

// Channels interacts with channels on the Insprd
func (c *Client) Channels() controller.ChannelInterface {
	return &ChannelClient{
		client: c.HTTPClient,
		config: c.Config,
	}
}

// Apps interacts with apps on the Insprd
func (c *Client) Apps() controller.AppInterface {
	return &AppClient{
		client: c.HTTPClient,
		config: c.Config,
	}
}

// ChannelTypes interacts with channel types on the Insprd
func (c *Client) ChannelTypes() controller.ChannelTypeInterface {
	return &ChannelTypeClient{
		client: c.HTTPClient,
		config: c.Config,
	}
}

// Authorization interacts with Insprd's auth
func (c *Client) Authorization() controller.AuthorizationInterface {
	return &AuthClient{
		client: c.HTTPClient,
		config: c.Config,
	}
}

// Alias interacts with alias on the Insprd
func (c *Client) Alias() controller.AliasInterface {
	return &AliasClient{
		client: c.HTTPClient,
		config: c.Config,
	}
}
