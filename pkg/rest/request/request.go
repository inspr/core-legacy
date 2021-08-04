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

// Send sends a request to the url specified in instantiation, with the given
// route and method, using
// the encoder to encode the body and the decoder to decode the response into
// the responsePtr
func (c Client) Send(ctx context.Context, route, method string, body, responsePtr interface{}) (err error) {
	buf, err := c.encoder(body)
	if err != nil {
		return ierrors.
			NewError().
			BadRequest().
			Message("error encoding body to json").
			InnerError(err).
			Build()
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.routeToURL(route),
		bytes.NewBuffer(buf),
	)
	if err != nil {
		return ierrors.
			NewError().
			BadRequest().
			Message("error creating request").
			InnerError(err).
			Build()
	}
	defer req.Body.Close()

	for key, values := range c.headers {
		req.Header[key] = values
	}

	if c.auth != nil {
		token, err := c.auth.GetToken()
		if err != nil {
			return ierrors.
				NewError().
				Unauthorized().
				Message("unable to get token from configuration").
				InnerError(err).
				Build()
		}
		req.Header.Add("Authorization", string(token))
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return ierrors.
			NewError().
			BadRequest().
			InnerError(err).
			Message(err.Error()).
			Build()
	}
	defer resp.Body.Close()

	err = c.handleResponseErr(resp)
	if err != nil {
		return err
	}

	updatedToken := resp.Header.Get("Authorization")
	if c.auth != nil && updatedToken != "" {
		err := c.auth.SetToken([]byte(updatedToken))
		if err != nil {
			return ierrors.
				NewError().
				BadRequest().
				Message("unable to update token").
				InnerError(err).
				Build()
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

func (c Client) handleResponseErr(resp *http.Response) error {
	decoder := c.decoderGenerator(resp.Body)
	var err *ierrors.InsprError
	defaultErr := ierrors.
		NewError().
		InternalServer().
		Message("cannot retrieve error from server").
		Build()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		decoder.Decode(&err)
		if err != nil {
			err.Wrap("status unauthorized")
			return err
		}
		return defaultErr
	case http.StatusForbidden:
		decoder.Decode(&err)
		if err != nil {
			err.Wrap("status forbidden")
			return err
		}
		return defaultErr
	case http.StatusNotFound:
		return ierrors.NewError().Message("route not found").Build()

	default:
		decoder.Decode(&err)
		if err == nil {
			return defaultErr
		}
		return err
	}
}
