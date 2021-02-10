package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	globalEnv "gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

const pollTimeout = 100

// Reader reads/commit messages from the channels defined in the env
type Reader struct {
	consumer    *kafka.Consumer
	lastMessage *kafka.Message
}

// NewReader return a new Reader
func NewReader() (*Reader, error) {
	kafkaEnv := GetEnvironment()
	globalEnv := globalEnv.GetEnvironment()

	var reader Reader

	newConsumer, errKafkaConsumer := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaEnv.KafkaBootstrapServers,
		"group.id":           globalEnv.InsprAppContext,
		"auto.offset.reset":  kafkaEnv.KafkaAutoOffsetReset,
		"enable.auto.commit": false,
	})
	if errKafkaConsumer != nil {
		return nil, ierrors.NewError().Message("failed to create a new kafka consumer").InnerError(errKafkaConsumer).InternalServer().Build()
	}
	reader.consumer = newConsumer

	channelsList := globalEnv.GetInputChannelList()
	if len(channelsList) == 0 {
		return nil, ierrors.NewError().Message("KAFKA_INPUT_CHANNELS not specified").InvalidChannel().Build()
	}

	channelsAsTopics := utils.Map(channelsList, toTopic)

	reader.consumer.SubscribeTopics(channelsAsTopics, nil)
	return &reader, nil
}

/*
ReadMessage reads message by message. Returns channel the message belongs to,
the message and an error if any occured.
*/
func (reader *Reader) ReadMessage() (string, interface{}, error) {
	for {
		event := reader.consumer.Poll(pollTimeout)
		switch ev := event.(type) {
		case *kafka.Message:

			channel := *ev.TopicPartition.Topic

			// Decoding Message
			message, errDecode := decode(ev.Value, fromTopic(channel).channel)
			if errDecode != nil {
				return "", nil, errDecode
			}

			reader.lastMessage = ev
			channelName := fromTopic(channel).channel
			return channelName, message, nil

		case kafka.Error:
			if ev.Code() == kafka.ErrAllBrokersDown {
				return "", nil, ierrors.
					NewError().
					InnerError(ev).
					Message("kafka error = all brokers are down" + ev.Error()).
					InternalServer().
					Build()
			}

		default:
			continue
		}
	}
}

// Commit commits the last message read by Reader
func (reader *Reader) Commit() error {
	_, errCommit := reader.consumer.CommitMessage(reader.lastMessage)
	if errCommit != nil {
		return ierrors.
			NewError().
			InnerError(errCommit).
			Message("failed to commit last message").
			InternalServer().
			Build()
	}
	return nil
}

// Close close the reader consumer
func (reader *Reader) Close() {
	reader.consumer.Close()
}
