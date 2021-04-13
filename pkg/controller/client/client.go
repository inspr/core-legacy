package client

import (
	"github.com/inspr/inspr/pkg/controller"
	"github.com/inspr/inspr/pkg/rest/request"
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
