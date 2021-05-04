package models

import "context"

// BROKER's interfaces

// Reader reads from a message broker
type Reader interface {
	ReadMessage(ctx context.Context, channel string) (BrokerData, error)
	Commit(ctx context.Context, channel string) error
}

// Writer writes messages in a message broker
type Writer interface {
	WriteMessage(channel string, msg interface{}) error
}
