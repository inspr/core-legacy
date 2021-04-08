package client

import (
	"gitlab.inspr.dev/inspr/core/pkg/controller"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

// Client implements communication with the Insprd
type Client struct {
	HTTPClient *request.Client
}

// NewControllerClient return a new Client
func NewControllerClient(rc *request.Client) controller.Interface {
	return &Client{
		HTTPClient: rc,
	}
}

// Channels interacts with channels on the Insprd
func (c *Client) Channels() controller.ChannelInterface {
	return &ChannelClient{
		c: c.HTTPClient,
	}
}

// Apps interacts with apps on the Insprd
func (c *Client) Apps() controller.AppInterface {
	return &AppClient{
		c: c.HTTPClient,
	}
}

// ChannelTypes interacts with channel types on the Insprd
func (c *Client) ChannelTypes() controller.ChannelTypeInterface {
	return &ChannelTypeClient{
		c: c.HTTPClient,
	}
}

// Authorization interacts with Insprd's auth
func (c *Client) Authorization() controller.AuthorizationInterface {
	return &AuthClient{
		c: c.HTTPClient,
	}
}
