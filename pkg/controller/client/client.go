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
func NewControllerClient(
	url string,
	auth request.Authenticator,
) controller.Interface {

	return &Client{
		HTTPClient: request.NewClient().
			BaseURL(url).
			Encoder(json.Marshal).
			Decoder(request.JSONDecoderGenerator).
			Authenticator(auth).
			Build(),
	}
}

// Channels interacts with channels on the Insprd
func (c *Client) Channels() controller.ChannelInterface {
	return &ChannelClient{
		rc: c.HTTPClient,
	}
}

// Apps interacts with apps on the Insprd
func (c *Client) Apps() controller.AppInterface {
	return &AppClient{
		rc: c.HTTPClient,
	}
}

// ChannelTypes interacts with channel types on the Insprd
func (c *Client) ChannelTypes() controller.ChannelTypeInterface {
	return &ChannelTypeClient{
		rc: c.HTTPClient,
	}
}

// Authorization interacts with Insprd's auth
func (c *Client) Authorization() controller.AuthorizationInterface {
	return &AuthClient{
		rc: c.HTTPClient,
	}
}

// Alias interacts with alias on the Insprd
func (c *Client) Alias() controller.AliasInterface {
	return &AliasClient{
		rc: c.HTTPClient,
	}
}
