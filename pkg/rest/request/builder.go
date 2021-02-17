package request

import (
	"encoding/json"
	"net/http"
)

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
	return NewClient().BaseURL(baseURL).Encoder(json.Marshal).Decoder(JSONDecoderGenerator).Build()
}
