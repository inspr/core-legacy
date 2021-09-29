package request

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/logs"
)

var logger *zap.Logger

func init() {
	logger, _ = logs.Logger(zap.Fields(zap.String("section", "sidecar-client"), zap.String("dapp-name", os.Getenv("INSPR_APP_ID"))))
}

const (
	// DefaultHost is the the standard hostname, it is used in requests made by
	// the cli to the insprd/uidp services in the cluster
	DefaultHost = "inspr.com"
)

var (
	// DefaultErr is the error returned by the request's response when an
	// unexpected http.Status is provided and it doesn't have an error structure
	// in its response body.
	DefaultErr = ierrors.
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
			ierrors.New(err).BadRequest(),
			"error encoding body to json",
		)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.routeToURL(route),
		bytes.NewBuffer(buf),
	)

	logger.Debug("Sending request to:" + c.routeToURL(route))

	if err != nil {
		return ierrors.Wrap(err, "error creating request")
	}

	for key, values := range c.headers {
		req.Header[key] = values
	}

	req.Host = c.host

	if c.auth != nil {
		token, err := c.auth.GetToken()
		if err != nil {
			return ierrors.Wrap(err, "unable to get token from configuration")
		}
		req.Header.Add("Authorization", string(token))
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return ierrors.New(err).BadRequest()
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
			return ierrors.Wrap(err, "unable to update token")
		}
	}

	if responsePtr != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(responsePtr)

		if errors.Is(err, io.EOF) {
			return nil
		}
	}

	return err
}

func (c Client) handleResponseErr(resp *http.Response) error {
	decoder := c.decoderGenerator(resp.Body)
	err := ierrors.New("")

	switch resp.StatusCode {
	case http.StatusOK:
		return nil

	case http.StatusUnauthorized:
		decoder.Decode(&err)
		return ierrors.Wrap(
			err,
			"status unauthorized",
		)
	case http.StatusForbidden:
		decoder.Decode(&err)
		return ierrors.Wrap(
			err,
			"status forbidden",
		)
	case http.StatusNotFound:
		return ierrors.New("route not found")

	default:
		decoder.Decode(&err)
		if err == nil {
			return DefaultErr
		}
		return err
	}
}
