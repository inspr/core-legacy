package client

import (
	"encoding/json"

	"github.com/inspr/inspr/pkg/controller"
	"github.com/inspr/inspr/pkg/rest/request"
)

// Client implements communication with the Insprd
type Client struct {
	HTTPClient *request.Client
}

// NewControllerClient return a new Client
func NewControllerClient(url string, auth request.Authenticator) controller.Interface {
	return &Client{
		HTTPClient: request.NewClient().
			BaseURL(url).
			Encoder(json.Marshal).
			Decoder(request.JSONDecoderGenerator).
			Authenticator(auth).
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
