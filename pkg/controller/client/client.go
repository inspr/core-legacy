package client

import (
	"gitlab.inspr.dev/inspr/core/pkg/controller"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

type Client struct {
	c *rest.Client
}

func (c *Client) Channels() controller.ChannelInterface {
	return &ChannelClient{
		c: c.c,
	}
}

func (c *Client) Apps() controller.AppInterface {
	return &AppClient{
		c: c.c,
	}
}

func (c *Client) ChannelTypes() controller.ChannelTypeInterface {
	return &ChannelTypeClient{
		c: c.c,
	}
}
