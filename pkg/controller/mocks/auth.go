package mocks

import (
	"context"

	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/controller"
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
func (am *AuthMock) GenerateToken(ctx context.Context, payload auth.Payload) (string, error) {
	return "", nil
}

// Init is the AuthMock Init method
func (am *AuthMock) Init(ctx context.Context, key string) (string, error) {
	return "", nil
}
