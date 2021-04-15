package dappclient

import (
	"context"

	"github.com/inspr/inspr/pkg/sidecar/models"
)

// AppClient defines an interface and its methods for an dApp Client
type AppClient interface {
	WriteMessage(ctx context.Context, channel string, msg models.Message) error
	ReadMessage(ctx context.Context, channel string, message interface{}) error
	CommitMessage(ctx context.Context, channel string) error
}
