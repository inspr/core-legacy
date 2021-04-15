// Package jwtauth is responsible for implementing the auth
// methods specified in the auth folder of the inspr pkg.
package jwtauth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/auth/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest/request"
	"github.com/lestrrat-go/jwx/jwt"
)

// JWTauth implements the Auth interface for jwt authetication provider
type JWTauth struct {
	authURL string
}

// NewJWTauth returns a JWTAuth object
func NewJWTauth() *JWTauth {
	url, ok := os.LookupEnv("AUTH_PATH")
	if !ok {
		panic("AUTH_PATH not found")
	}
	return &JWTauth{
		authURL: url,
	}
}

//Validate receives the
func (JA *JWTauth) Validate(token []byte) (models.Payload, []byte, error) {
	jwtToken, err := jwt.Parse(token, jwt.WithValidate(true))
	// not valid token
	if err != nil {
		return models.Payload{}, token, err
	}

	expiration, found := jwtToken.Get(jwt.ExpirationKey)
	now := time.Now()

	// expired
	if now.After(expiration.(time.Time)) || !found {
		newToken, err := JA.Refresh(token)
		if err != nil {
			return models.Payload{},
				token,
				errors.New("error refreshing token")
		}
		token = newToken
	}

	// gets payload from token
	payload, found := jwtToken.Get("payload")
	if !found {
		return models.Payload{},
			token,
			errors.New("payload not found")
	}

	return payload.(models.Payload), token, nil
}

// Tokenize receives a payload and returns it in signed jwt format. Uses JWT authentication provider
func (JA *JWTauth) Tokenize(load models.Payload) ([]byte, error) {

	client := request.NewJSONClient(JA.authURL)

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
	client := request.NewJSONClient(JA.authURL)

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
