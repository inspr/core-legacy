package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	globalEnv "gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

type readerObj struct {
	consumer    *kafka.Consumer
	lastMessage *kafka.Message
}

// Reader reads messages from the channels defined in the env
type Reader interface {
	Commit() error
	ReadMessage() (string, interface{}, error)
	Close()
}

// NewReader return a new reader
func NewReader() (Reader, error) {
	kafkaEnv := GetEnvironment()
	globalEnv := globalEnv.GetEnvironment()

	var reader readerObj

	newConsumer, errKafkaConsumer := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaEnv.KafkaBootstrapServers,
		"group.id":           globalEnv.InsprNodeID,
		"auto.offset.reset":  kafkaEnv.KafkaAutoOffsetReset,
		"enable.auto.commit": false,
	})
	if errKafkaConsumer != nil {
		return nil, ierrors.NewError().Build()
	}
	reader.consumer = newConsumer

	channelsList := globalEnv.InputChannels
	if len(channelsList) == 0 {
		return nil, ierrors.NewError().Build()
	}

}

// Commit
func (r *readerObj) Commit() error {
	_, errCommit := r.consumer.CommitMessage(r.lastMessage)
	if errCommit != nil {
		return ierrors.NewError().Build()
	}
	return nil
}
