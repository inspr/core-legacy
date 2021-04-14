// package jwtauth is responsible for implementing the auth
// methods specified in the auth folder of the inspr pkg.
package jwtauth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/jwt"
	"gitlab.inspr.dev/inspr/core/pkg/auth"
	"gitlab.inspr.dev/inspr/core/pkg/auth/models"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

type JWTauth struct{}

func NewJWTauth() *JWTauth {
	return &JWTauth{}
}

// receives the
func (JA *JWTauth) Validade(token []byte) (models.Payload, []byte, error) {
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
