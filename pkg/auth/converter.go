package auth

import (
	"encoding/json"

	"github.com/lestrrat-go/jwx/jwt"
	"gitlab.inspr.dev/inspr/core/pkg/auth/models"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// Desserialize converts a interface to a Payload model
func Desserialize(tokenBytes []byte) (*models.Payload, error) {

	token, err := jwt.Parse(tokenBytes)
	if err != nil {
		err = ierrors.NewError().InternalServer().Message("error: didn't return a token").Build()
		return nil, err
	}
	load, ok := token.Get("payload")
	if !ok {
		err = ierrors.NewError().InternalServer().Message("error: didn't return a payload on it's token").Build()
		return nil, err
	}

	jwtJSON, err := json.Marshal(load)
	if err != nil {
		return nil, err
	}

	var payload models.Payload
	err = json.Unmarshal(jwtJSON, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
