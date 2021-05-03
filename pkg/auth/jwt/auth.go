// Package jwtauth is responsible for implementing the auth
// methods specified in the auth folder of the inspr pkg.
package jwtauth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest/request"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

// JWTauth implements the Auth interface for jwt authetication provider
type JWTauth struct {
	publicKey *rsa.PublicKey
	authURL   string
}

// NewJWTauth takes an *rsa.PublicKey and returns an
// structure that implements the auth interface
func NewJWTauth(rsaPublicKey *rsa.PublicKey) *JWTauth {
	url, ok := os.LookupEnv("AUTH_PATH")
	if !ok {
		panic("[ENV VAR] AUTH_PATH not found")
	}
	return &JWTauth{
		publicKey: rsaPublicKey,
		authURL:   url,
	}
}

// Validate is a wrapper that checks the token of the http request and if it's
// valid, proceeds to execute the request and if it isn't valid returns an error
func (JA *JWTauth) Validate(token []byte) (*auth.Payload, []byte, error) {

	_, err := jwt.Parse(
		token,
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, JA.publicKey),
	)
	if err != nil {
		if err.Error() == errors.New(`exp not satisfied`).Error() {

			newToken, err := JA.Refresh(token)
			if err != nil {
				return nil,
					token,
					ierrors.
						NewError().
						InternalServer().
						InnerError(err).
						Message("error refreshing token").
						Build()
			}
			token = newToken
		} else {
			return nil, token, err
		}
	}

	// expired

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

// InitDO  structure for initialization requests
type InitDO struct {
	auth.Payload
	Key string
}

// Init receives a payload and returns it in signed jwt format. Uses JWT authentication provider
func (JA *JWTauth) Init(key string, load auth.Payload) ([]byte, error) {
	initDO := InitDO{
		Key:     key,
		Payload: load,
	}

	client := request.NewJSONClient(JA.authURL)

	data := auth.JwtDO{}
	err := client.Send(context.Background(), "/init", http.MethodPost, initDO, &data)
	if err != nil {
		log.Printf("err = %+v\n", err)
		err = ierrors.NewError().InternalServer().Message(err.Error()).Build()
		return nil, err
	}

	return data.Token, nil
}

// Tokenize receives a payload and returns it in signed jwt format. Uses JWT authentication provider
func (JA *JWTauth) Tokenize(load auth.Payload) ([]byte, error) {

	client := request.NewJSONClient(JA.authURL)

	data := auth.JwtDO{}
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

	log.Printf("string(token) = %+v\n", string(token))
	data := auth.JwtDO{}

	err := client.Send(context.Background(), "/refresh", http.MethodGet, nil, &data)
	if err != nil {
		err = ierrors.NewError().InternalServer().Message(err.Error()).Build()
		return nil, err
	}
	log.Printf("string(data.Token) = %+v\n", string(data.Token))

	return data.Token, nil
}
