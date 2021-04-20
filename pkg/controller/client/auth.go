package client

import (
	"context"

	"github.com/inspr/inspr/pkg/api/auth"
	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/rest/request"
)

// AuthClient is a client for getting auth information from Insprd
type AuthClient struct {
	c *request.Client
}

// GenerateToken sends a request containing a payload so Insprd
// generates a new auth token based on the payload's info
func (ac *AuthClient) GenerateToken(ctx context.Context, payload auth.Payload) (string, error) {
	authDI := models.AuthDI{}

	err := ac.c.Send(ctx, "/auth", "POST", payload, &authDI)
	if err != nil {
		return "", err
	}

	return authDI.Token, nil
}