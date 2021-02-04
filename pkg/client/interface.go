package dappclient

import (
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

// AppClient todo doc
type AppClient interface {
	WriteMessage(channel string, msg models.Message) error
	ReadMessage(channel string) (models.Message, error)
	CommitMessage(channel string) error
}
