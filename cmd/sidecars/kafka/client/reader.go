package kafkasc

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	globalEnv "gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

const pollTimeout = 100

// Consumer interface
type Consumer interface {
	Poll(timeout int) (event kafka.Event)
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) (err error)
	CommitMessage(m *kafka.Message) ([]kafka.TopicPartition, error)
	Close() (err error)
}

// Reader reads/commit messages from the channels defined in the env
type Reader struct {
	consumer    Consumer
	lastMessage *kafka.Message
}

// NewReader return a new Reader
func NewReader() (*Reader, error) {
	kafkaEnv := GetEnvironment()

	var reader Reader

	newConsumer, errKafkaConsumer := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaEnv.KafkaBootstrapServers,
		"group.id":           globalEnv.GetInsprAppID(),
		"auto.offset.reset":  kafkaEnv.KafkaAutoOffsetReset,
		"enable.auto.commit": false,
	})

	if errKafkaConsumer != nil {
		return nil, ierrors.NewError().Message(errKafkaConsumer.Error()).InnerError(errKafkaConsumer).InternalServer().Build()
	}
	reader.consumer = newConsumer

	channelsList := globalEnv.GetResolvedInputChannelList(globalEnv.GetInputChannels())
	if len(channelsList) == 0 {
		return nil, ierrors.NewError().Message("KAFKA_INPUT_CHANNELS not specified").InvalidChannel().Build()
	}

	channelsAsTopics := channelsList.Map(func(s string) string {
		ch, _ := fromResolvedChannel(s)
		return ch.toTopic()
	})

	if err := reader.consumer.SubscribeTopics(channelsAsTopics, nil); err != nil {
		return nil, err
	}
	return &reader, nil
}

/*
ReadMessage reads message by message. Returns channel the message belongs to,
the message and an error if any occurred.
*/
func (reader *Reader) ReadMessage() (models.BrokerData, error) {
	for {
		event := reader.consumer.Poll(pollTimeout)
		switch ev := event.(type) {
		case *kafka.Message:

			topic := *ev.TopicPartition.Topic
			channel := fromTopic(topic)

			// Decoding Message
			message, errDecode := channel.decode(ev.Value)
			if errDecode != nil {
				return models.BrokerData{}, errDecode
			}

			reader.lastMessage = ev
			channelName := channel.channel

			return models.BrokerData{Message: models.Message{Data: message}, Channel: channelName}, nil

		case kafka.Error:
			if ev.Code() == kafka.ErrAllBrokersDown {
				return models.BrokerData{}, ierrors.
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

// CommitMessage commits the last message read by Reader
func (reader *Reader) CommitMessage() error {
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
func (reader *Reader) Close() error {
	err := reader.consumer.Close()
	if err != nil {
		return err
	}
	return nil
}
