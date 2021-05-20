package models

import "context"

// Reader reads from a message broker
type Reader interface {
	ReadMessage(ctx context.Context, channel string) (BrokerMessage, error)
	Commit(ctx context.Context, channel string) error
}

// Writer writes messages in a message broker
type Writer interface {
	WriteMessage(channel string, msg interface{}) error
}
