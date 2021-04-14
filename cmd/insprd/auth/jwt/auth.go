// package jwtauth is responsible for implementing the auth
// methods specified in the auth folder of the inspr pkg.
package jwtauth

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"gitlab.inspr.dev/inspr/core/pkg/auth/models"
)

type JWTauth struct {
	rsaKey *rsa.PrivateKey
}

func NewJWTauth(privateKey *rsa.PrivateKey) *JWTauth {
	return &JWTauth{
		rsaKey: privateKey,
	}
}

// Validade is a wrapper that checks the token of the http request and if it's
// valid, proceeds to execute the request and if it isn't valid returns an error
func (JA *JWTauth) Validade(token []byte) (models.Payload, []byte, error) {

	jwtToken, err := jwt.Parse(
		token,
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, &JA.rsaKey.PublicKey),
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

func (JA *JWTauth) Tokenize(load models.Payload) ([]byte, error) {
	return []byte{}, nil
}
func (JA *JWTauth) Refresh(token []byte) ([]byte, error) {
	return []byte{}, nil
}
