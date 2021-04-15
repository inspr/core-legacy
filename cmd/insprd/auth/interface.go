package auth

import "github.com/inspr/inspr/pkg/auth/models"

//Auth is the inteface for interacting with the Authentication service
type Auth interface {
	Validate(token []byte) (models.Payload, []byte, error)
	Tokenize(load models.Payload) ([]byte, error)
	Refresh(token []byte) ([]byte, error)
}
