package auth

import (
	"encoding/json"

	"github.com/inspr/inspr/pkg/auth/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/lestrrat-go/jwx/jwt"
)

// Desserialize converts a interface to a Payload model
func Desserialize(tokenBytes []byte) (*models.Payload, error) {

	token, err := jwt.Parse(tokenBytes)
	if err != nil {
		err = ierrors.NewError().InternalServer().Message("jwt parsing failed, error: %s", err.Error()).Build()
		return nil, err
	}

	load, ok := token.Get("payload")
	if !ok {
		err = ierrors.NewError().InternalServer().Message("error: jwt token didn't carry a payload").Build()
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
