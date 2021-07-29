package request

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"inspr.dev/inspr/pkg/ierrors"
)

const (
	// DefaultHost is the the standard hostname, it is used in requests made by
	// the cli to the insprd/uidp services in the cluster
	DefaultHost = "inspr.com"
)

var (
	defaultErr = ierrors.
		New("cannot retrieve error from server").
		InternalServer()
)

// Send sends a request to the url specified in instantiation, with the given
// route and method, using
// the encoder to encode the body and the decoder to decode the response into
// the responsePtr
func (c Client) Send(ctx context.Context, route, method string, body, responsePtr interface{}) (err error) {
	buf, err := c.encoder(body)
	if err != nil {
		return ierrors.Wrap(
			ierrors.From(err).BadRequest(),
			"error encoding body to json",
		)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.routeToURL(route),
		bytes.NewBuffer(buf),
	)

	if err != nil {
		return ierrors.Wrap(
			ierrors.From(err).BadRequest(),
			"error creating request",
		)
	}

	for key, values := range c.headers {
		req.Header[key] = values
	}

	if c.auth != nil {
		token, err := c.auth.GetToken()
		if err != nil {
			return ierrors.Wrap(
				ierrors.From(err).BadRequest(),
				"unable to get token from configuration",
			)
		}
		req.Header.Add("Authorization", string(token))
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return ierrors.From(err).BadRequest()
	}

	err = c.handleResponseErr(resp)
	if err != nil {
		return err
	}

	updatedToken := resp.Header.Get("Authorization")
	if c.auth != nil && updatedToken != "" {
		err := c.auth.SetToken([]byte(updatedToken))
		if err != nil {
			return ierrors.Wrap(
				ierrors.From(err).BadRequest(),
				"unable to update token",
			)
		}
	}

	if responsePtr != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(responsePtr)

		if err == io.EOF {
			return nil
		}

	}

	return err
}

// TODO REVIEW

func (c Client) handleResponseErr(resp *http.Response) error {
	decoder := c.decoderGenerator(resp.Body)
	err := ierrors.New("")

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		decoder.Decode(&err)
		if err != nil {
			ierrors.Wrap(
				err,
				"status unauthorized",
			)
			return err
		}
		return defaultErr
	case http.StatusForbidden:
		decoder.Decode(&err)
		if err != nil {
			ierrors.Wrap(
				err,
				"status forbidden",
			)
			return err
		}
		return defaultErr
	case http.StatusNotFound:
		return ierrors.New("route not found")

	default:
		decoder.Decode(&err)
		if err == nil {
			return defaultErr
		}
		return err
	}
}
