package auth

import (
	"encoding/json"

	"github.com/lestrrat-go/jwx/jwt"
	"inspr.dev/inspr/pkg/ierrors"
)

// Desserialize converts a token to a Payload model
func Desserialize(tokenBytes []byte) (*Payload, error) {

	token, err := jwt.Parse(tokenBytes)
	if err != nil {
		err = ierrors.Wrap(
			ierrors.New(err).InternalServer(),
			"jwt parsing failed",
		)
		return nil, err
	}

	load, ok := token.Get("payload")
	if !ok {
		err = ierrors.New(
			"jwt token didn't carry a payload",
		).InternalServer()
		return nil, err
	}

	jwtJSON, err := json.Marshal(load)
	if err != nil {
		return nil, ierrors.New(err).InternalServer()
	}

	var payload Payload
	err = json.Unmarshal(jwtJSON, &payload)
	if err != nil {
		return nil, ierrors.New(err).InternalServer()
	}
	return &payload, nil
}
