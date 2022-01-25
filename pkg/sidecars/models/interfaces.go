package models

import (
	"context"
)

// Reader reads from a message broker
type Reader interface {
	ReadMessage(ctx context.Context, channel string) ([]byte, error)
	Commit(ctx context.Context, channel string) error
	Close() error
}

// Writer writes messages in a message broker
type Writer interface {
	WriteMessage(channel string, msg []byte) error
	Close()
}

type BrokerInterface interface {
	Reader() Reader
	Writer() Writer
}
