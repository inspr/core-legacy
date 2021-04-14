package client

import (
	"context"

	"github.com/inspr/inspr/cmd/insprd/auth"
	"github.com/inspr/inspr/pkg/rest/request"
)

// AuthClient is a client for getting auth information from Insprd
type AuthClient struct {
	c *request.Client
}

// GenerateToken sends a request containing a payload so Insprd
// generates a new auth token based on the payload's info
func (ac *AuthClient) GenerateToken(ctx context.Context, payload auth.Payload) (string, error) {
	type returnStructure struct {
		Token string `json:"token"`
	}
	var resp returnStructure

	err := ac.c.Send(ctx, "/auth", "POST", payload, &resp)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}
