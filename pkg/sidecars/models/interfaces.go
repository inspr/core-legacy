package models

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer interface
type Consumer interface {
	Poll(int) kafka.Event
	Commit() ([]kafka.TopicPartition, error)
	Close() (err error)
}

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
