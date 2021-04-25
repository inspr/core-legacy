package mocks

import (
	"context"

	"github.com/inspr/inspr/pkg/api/auth"
	"github.com/inspr/inspr/pkg/controller"
)

// AuthMock mock structure for the operations of the controller.Authorization()
type AuthMock struct {
	err error
}

// NewAuthMock exports a mock of the Authorization.interface
func NewAuthMock(err error) controller.AuthorizationInterface {
	return &AuthMock{err: err}
}

// GenerateToken is the AuthMock GenerateToken method
func (ac *AuthMock) GenerateToken(ctx context.Context, payload auth.Payload) (string, error) {
	return "", nil
}
func (ac *AuthMock) Init(ctx context.Context, key string) (string, error) {
	return "", nil
}
