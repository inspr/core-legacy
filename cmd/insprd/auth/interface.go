package auth

import (
	"gitlab.inspr.dev/inspr/core/pkg/auth"
)

//Auth is the inteface for interacting with the Authentication service
type Auth interface {
	Validade(token []byte) ([]byte, error)
	Tokenize(load auth.Payload) ([]byte, error)
	Refresh(token []byte) ([]byte, error)
}
