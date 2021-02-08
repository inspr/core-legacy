package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// Send sends a request to the url specified in instantiation, with the given route and method, using
// the encoder to encode the body and the decoder to decode the response into the responsePtr
func (c *Client) Send(ctx context.Context, route string, method string, body interface{}, responsePtr interface{}) (err error) {
	buf, err := c.encoder(body)
	if err != nil {
		return ierrors.NewError().BadRequest().Message("error encoding body to json").InnerError(err).Build()
	}

	req, err := http.NewRequest(method, c.routeToURL(route), bytes.NewBuffer(buf))
	if err != nil {
		return ierrors.NewError().BadRequest().Message("error creating request").InnerError(err).Build()
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return ierrors.NewError().BadRequest().InnerError(err).Message("unable to send request to insprd").Build()
	}

	err = c.handleResponseErr(resp)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(responsePtr)

	return err
}

func (c *Client) routeToURL(route string) string {
	return fmt.Sprintf("%s%s", c.baseURL, route)
}

// Encoder encodes an interface into bytes
type Encoder func(interface{}) ([]byte, error)

// DecoderGenerator creates a decoder for a given request
type DecoderGenerator func(r io.Reader) Decoder

// JSONDecoderGenerator generates a decoder for json encoded requests
func JSONDecoderGenerator(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

// Decoder is an interface that decodes a reader into an struct
type Decoder interface {
	Decode(interface{}) error
}

// Client is a generic rest client
type Client struct {
	c                http.Client
	baseURL          string
	encoder          Encoder
	decoderGenerator DecoderGenerator
}

func (c *Client) handleResponseErr(resp *http.Response) error {
	decoder := c.decoderGenerator(resp.Body)
	var err *ierrors.InsprError
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		decoder.Decode(&err)
		return err
	}
}
