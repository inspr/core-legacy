// Package jwtauth is responsible for implementing the auth
// methods specified in the auth folder of the inspr pkg.
package jwtauth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
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

// JWTauth implements the Auth interface for jwt authetication provider
type JWTauth struct {
	PublicKey *rsa.PublicKey
	authURL   string
}

// NewJWTauth takes an *rsa.PublicKey and returns an
// structure that implements the auth interface
func NewJWTauth(rsaPublicKey *rsa.PublicKey) *JWTauth {
	url, ok := os.LookupEnv("AUTH_PATH")
	if !ok {
		panic("AUTH_PATH not found")
	}
	return &JWTauth{
		PublicKey: rsaPublicKey,
		authURL:   url,
	}
}

// Validate is a wrapper that checks the token of the http request and if it's
// valid, proceeds to execute the request and if it isn't valid returns an error
func (JA *JWTauth) Validate(token []byte) (*models.Payload, []byte, error) {

	jwtToken, err := jwt.Parse(
		token,
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, JA.PublicKey),
	)

	// not valid token
	if err != nil {
		return nil, token, err
	}

	expiration, found := jwtToken.Get(jwt.ExpirationKey)
	now := time.Now()

	// expired
	if !found || now.After(expiration.(time.Time)) {
		newToken, err := JA.Refresh(token)
		if err != nil {
			return nil,
				token,
				ierrors.
					NewError().
					InternalServer().
					Message("error refreshing token").
					Build()
		}
		token = newToken
	}

	// gets payload from token
	payload, err := auth.Desserialize(token)
	if err != nil {
		return nil,
			token,
			ierrors.
				NewError().
				InternalServer().
				Message("error desserializing the payload").
				Build()
	}

	return payload, token, nil
}

// Tokenize receives a payload and returns it in signed jwt format. Uses JWT authentication provider
func (JA *JWTauth) Tokenize(load models.Payload) ([]byte, error) {

	client := request.NewJSONClient(JA.authURL)

	data := models.JwtDO{}
	err := client.Send(context.Background(), "/token", http.MethodPost, load, &data)
	if err != nil {
		err = ierrors.NewError().InternalServer().Message(err.Error()).Build()
		return nil, err
	}

	return data.Token, nil
}

// Refresh refreshes a jwt token. Uses JWT authentication provider
func (JA *JWTauth) Refresh(token []byte) ([]byte, error) {
	client := request.NewClient().
		BaseURL(JA.authURL).
		Encoder(json.Marshal).
		Decoder(request.JSONDecoderGenerator).
		Header("Authorization", fmt.Sprintf("Bearer %v", string(token))).
		Build()

	data := models.JwtDO{}

	err := client.Send(context.Background(), "/refresh", http.MethodGet, nil, &data)
	if err != nil {
		err = ierrors.NewError().InternalServer().Message(err.Error()).Build()
		return nil, err
	}

	return data.Token, nil
}
