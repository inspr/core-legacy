package auth

import "gitlab.inspr.dev/inspr/core/pkg/auth/models"

//Auth is the inteface for interacting with the Authentication service
type Auth interface {
	Validade(token []byte) (models.Payload, []byte, error)
	Tokenize(load models.Payload) ([]byte, error)
	Refresh(token []byte) ([]byte, error)
}