package client

import (
	"context"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/auth"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

// AuthClient is a client for getting auth information from Insprd
type AuthClient struct {
	c *request.Client
}

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
