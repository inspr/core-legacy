package request

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: move this to another package, it doens't really fit the rest pkg.

// Authenticator is an interface to perform authentication via tokens
type Authenticator interface {
	GetToken() ([]byte, error)
	SetToken([]byte) error
}

// NewClient returns an address of a empty Client
func NewClient() *Client {
	return &Client{}
}

// BaseURL sets the base URL for the client that is being built
func (c *Client) BaseURL(url string) *Client {
	c.baseURL = url
	return c
}

// Encoder sets the encoder for the client that is being built
func (c *Client) Encoder(encoder Encoder) *Client {
	c.encoder = encoder
	return c
}

// Decoder sets the decoder for the client that is being built
func (c *Client) Decoder(decoder DecoderGenerator) *Client {
	c.decoderGenerator = decoder
	return c
}

// Authenticator adds the authentication interface implementation to the
// Client strucuture.
func (c *Client) Authenticator(au Authenticator) *Client {
	c.auth = au
	return c
}

// Token adds a token header with the format "Authentication: Bearer " + token on each request the client sends.
func (c *Client) Token(token []byte) *Client {
	return c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
}

// HTTPClient sets the http client for the client that is being built
func (c *Client) HTTPClient(client http.Client) *Client {
	c.c = client
	return c
}

// Build returns the client built by the builder
func (c *Client) Build() *Client {
	return c
}

// Header adds the given header to all requests made by the client
func (c *Client) Header(key, value string) *Client {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	c.headers[key] = value
	return c
}

// NewJSONClient returns a client for the given url with json encoding and decoding
func NewJSONClient(baseURL string) *Client {
	client := &Client{}
	return client.
		BaseURL(baseURL).
		Encoder(json.Marshal).
		Decoder(JSONDecoderGenerator).
		Build()
}
