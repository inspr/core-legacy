package auth

import (
	"encoding/json"

	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/lestrrat-go/jwx/jwt"
)

// Desserialize converts a interface to a Payload model
func Desserialize(tokenBytes []byte) (*Payload, error) {

	token, err := jwt.Parse(tokenBytes)
	if err != nil {
		err = ierrors.NewError().InternalServer().Message("jwt parsing failed, error: %s", err.Error()).Build()
		return nil, err
	}

	load, ok := token.Get("payload")
	if !ok {
		err = ierrors.NewError().InternalServer().Message("jwt token didn't carry a payload").Build()
		return nil, err
	}

	jwtJSON, err := json.Marshal(load)
	if err != nil {
		return nil, err
	}

	var payload Payload
	err = json.Unmarshal(jwtJSON, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
