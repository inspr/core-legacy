package request

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Authenticator is an interface to perform authentication via tokens
type Authenticator interface {
	GetToken() ([]byte, error)
	SetToken([]byte) error
}

// ClientBuilder builds a client with the given specifications
type ClientBuilder struct {
	c *Client
}

// BaseURL sets the base URL for the client that is being built
func (cb *ClientBuilder) BaseURL(url string) *ClientBuilder {
	cb.c.baseURL = url
	return cb
}

// Encoder sets the encoder for the client that is being built
func (cb *ClientBuilder) Encoder(encoder Encoder) *ClientBuilder {
	cb.c.encoder = encoder
	return cb
}

// Decoder sets the decoder for the client that is being built
func (cb *ClientBuilder) Decoder(decoder DecoderGenerator) *ClientBuilder {
	cb.c.decoderGenerator = decoder
	return cb
}

// Authenticator adds the authentication interface implementation to the
// Client strucuture.
func (cb *ClientBuilder) Authenticator(au Authenticator) *ClientBuilder {
	cb.c.auth = au
	return cb
}

// Token adds a token header with the format "Authentication: Bearer " + token on each request the client sends.
func (cb *ClientBuilder) Token(token []byte) *ClientBuilder {
	return cb.Header("Authorization", fmt.Sprintf("Bearer %s", token))
}

// HTTPClient sets the http client for the client that is being built
func (cb *ClientBuilder) HTTPClient(client http.Client) *ClientBuilder {
	cb.c.c = client
	return cb
}

// NewClient creates a builder for a client
func NewClient() *ClientBuilder {
	return &ClientBuilder{
		c: &Client{},
	}
}

// Build returns the client built by the builder
func (cb *ClientBuilder) Build() *Client {
	return cb.c
}

// NewJSONClient returns a client for the given url with json encoding and decoding
func NewJSONClient(baseURL string) *Client {
	return NewClient().
		BaseURL(baseURL).
		Encoder(json.Marshal).
		Decoder(JSONDecoderGenerator).
		Build()
}

// Header adds the given header to all requests made by the client
func (cb *ClientBuilder) Header(key, value string) *ClientBuilder {
	if cb.c.headers == nil {
		cb.c.headers = make(map[string]string)
	}
	cb.c.headers[key] = value
	return cb
}
