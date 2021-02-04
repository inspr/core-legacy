package dappclient

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

// AppClient todo doc
type AppClient interface {
	WriteMessage(ctx context.Context, channel string, msg models.Message) error
	ReadMessage(ctx context.Context, channel string) (models.Message, error)
	CommitMessage(ctx context.Context, channel string) error
}
