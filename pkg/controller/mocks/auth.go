package mocks

import (
	"context"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/auth"
	"gitlab.inspr.dev/inspr/core/pkg/controller"
)

// AuthMock mock structure for the operations of the controller.Authorization()
type AuthMock struct {
	err error
}

// NewAuthMock exports a mock of the Authorization.interface
func NewAuthMock(err error) controller.AuthorizationInterface {
	return &AuthMock{err: err}
}

func (ac *AuthMock) GenerateToken(ctx context.Context, payload auth.Payload) (string, error) {
	return "", nil
}
