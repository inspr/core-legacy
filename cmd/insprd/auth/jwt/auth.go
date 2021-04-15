// Package jwtauth is responsible for implementing the auth
// methods specified in the auth folder of the inspr pkg.
package jwtauth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/auth/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest/request"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

// JWTauth structure containing the private key of in the service side,
// this key is used to parse the user keys given in the requests.
type JWTauth struct {
	PublicKey *rsa.PublicKey
}

// NewJWTauth takes an *rsa.PrivateKey and returns an
// structure that implements the auth interface
func NewJWTauth(rsaPublicKey *rsa.PublicKey) *JWTauth {
	return &JWTauth{
		PublicKey: rsaPublicKey,
	}
}

// Validade is a wrapper that checks the token of the http request and if it's
// valid, proceeds to execute the request and if it isn't valid returns an error
func (JA *JWTauth) Validate(token []byte) (models.Payload, []byte, error) {

	jwtToken, err := jwt.Parse(
		token,
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, JA.PublicKey),
	)

	// not valid token
	if err != nil {
		return models.Payload{}, token, err
	}

	expiration, found := jwtToken.Get(jwt.ExpirationKey)
	now := time.Now()

	// expired
	if !found || now.After(expiration.(time.Time)) {
		newToken, err := JA.Refresh(token)
		if err != nil {
			return models.Payload{},
				token,
				errors.New("error refreshing token")
		}
		token = newToken
	}

	// gets payload from token
	payloadData, found := jwtToken.Get("payload")
	if !found {
		return models.Payload{},
			token,
			errors.New("payload not found")
	}

	// gets the string in the payload and converts it to bytes
	payloadString := payloadData.(string)
	payloadBytes := []byte(payloadString)

	// unmarshal of the bytes
	var payload models.Payload
	json.Unmarshal(payloadBytes, &payload)

	return payload, token, nil
}

// Tokenize receives a payload and returns it in signed jwt format. Uses JWT authentication provider
func (JA *JWTauth) Tokenize(load models.Payload) ([]byte, error) {

	URL := os.Getenv("AUTH_PATH")
	client := request.NewJSONClient(URL)

	data := models.JwtDO{}
	err := client.Send(context.Background(), "/token", http.MethodPost, load, &data)
	if err != nil {
		err = ierrors.NewError().InternalServer().InnerError(err).Build()
		return nil, err
	}

	return data.Token, nil
}

// Refresh refreshes a jwt token. Uses JWT authentication provider
func (JA *JWTauth) Refresh(token []byte) ([]byte, error) {
	URL := os.Getenv("AUTH_PATH")
	client := request.NewJSONClient(URL)

	load, err := auth.Desserialize(token)
	if err != nil {
		return nil, err
	}

	body := models.ResfreshDI{
		RefreshToken: load.Refresh,
		RefreshURL:   load.RefreshURL,
	}

	data := models.JwtDO{}
	err = client.Send(context.Background(), "/refresh", http.MethodPost, body, &data)

	if err != nil {
		err = ierrors.NewError().InternalServer().InnerError(err).Build()
		return nil, err
	}

	return data.Token, nil
}
