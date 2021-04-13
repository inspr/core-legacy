// package jwtauth is responsible for implementing the auth
// methods specified in the auth folder of the inspr pkg.
package jwtauth

import "gitlab.inspr.dev/inspr/core/pkg/auth/models"

type JWTauth struct{}

func (JA *JWTauth) Validade(token []byte) (models.Payload, []byte, error) {

	return models.Payload{}, []byte{}, nil
}

func (JA *JWTauth) Tokenize(load models.Payload) ([]byte, error)
func (JA *JWTauth) Refresh(token []byte) ([]byte, error)
