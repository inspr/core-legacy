package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Encoder encodes an interface into bytes
type Encoder func(interface{}) ([]byte, error)

// DecoderGenerator creates a decoder for a given request
type DecoderGenerator func(r io.Reader) Decoder

// Decoder is an interface that decodes a reader into an struct
type Decoder interface{ Decode(interface{}) error }

// JSONDecoderGenerator generates a decoder for json encoded requests
func JSONDecoderGenerator(r io.Reader) Decoder { return json.NewDecoder(r) }

type stringSlice []string

// Client is a generic rest client
type Client struct {
	c                http.Client
	baseURL          string
	encoder          Encoder
	decoderGenerator DecoderGenerator
	headers          map[string]stringSlice
	auth             Authenticator
}

func (c Client) routeToURL(route string) string {
	return fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(route, "/"))
}

// Authenticator is an interface to perform authentication via tokens
type Authenticator interface {
	GetToken() ([]byte, error)
	SetToken([]byte) error
}

// NewClient returns an address of a empty Client
func NewClient() Client {
	return Client{}
}

// NewJSONClient returns a client for the given url with json encoding and decoding
func NewJSONClient(baseURL string) *Client {
	return NewClient().
		BaseURL(baseURL).
		Encoder(json.Marshal).
		Decoder(JSONDecoderGenerator).
		Build()
}

// BaseURL sets the base URL for the client that is being built
func (c Client) BaseURL(url string) Client {
	c.baseURL = url
	return c
}

// Build returns the address of an existing client
func (c Client) Build() *Client {
	return &c
}

// Encoder sets the encoder for the client that is being built
func (c Client) Encoder(encoder Encoder) Client {
	c.encoder = encoder
	return c
}

// Decoder sets the decoder for the client that is being built
func (c Client) Decoder(decoder DecoderGenerator) Client {
	c.decoderGenerator = decoder
	return c
}

// Authenticator adds the authentication interface implementation to the
// Client strucuture.
func (c Client) Authenticator(au Authenticator) Client {
	c.auth = au
	return c
}

// Token adds a token header with the format "Authentication: Bearer " + token on each request the client sends.
func (c Client) Token(token []byte) Client {
	return c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
}

// HTTPClient sets the http client for the client that is being built
func (c Client) HTTPClient(client http.Client) Client {
	c.c = client
	return c
}

// Header adds the value into the slice located by the key in the client's header map.
func (c Client) Header(key, value string) Client {
	if c.headers == nil {
		c.headers = make(map[string]stringSlice)
	}
	c.headers[key] = append(c.headers[key], value)
	return c
}
