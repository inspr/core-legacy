package kafkasc

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	globalEnv "gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
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
	// consumer    Consumer //deprecated
	consumers   map[string]Consumer
	lastMessage *kafka.Message
}

// NewReader return a new Reader
func NewReader() (*Reader, error) {
	// kafkaEnv := GetEnvironment()

	var reader Reader

	// newConsumer, errKafkaConsumer := kafka.NewConsumer(&kafka.ConfigMap{
	// 	"bootstrap.servers":  kafkaEnv.KafkaBootstrapServers,
	// 	"group.id":           globalEnv.GetInsprAppID(),
	// 	"auto.offset.reset":  kafkaEnv.KafkaAutoOffsetReset,
	// 	"enable.auto.commit": false,
	// })

	// if errKafkaConsumer != nil {
	// 	return nil, ierrors.NewError().Message(errKafkaConsumer.Error()).InnerError(errKafkaConsumer).InternalServer().Build()
	// }
	// reader.consumer = newConsumer

	channelsList := globalEnv.GetInputChannelList(globalEnv.GetInputChannels())
	if len(channelsList) == 0 {
		return nil, ierrors.NewError().Message("KAFKA_INPUT_CHANNELS not specified").InvalidChannel().Build()
	}

	channelsAsTopics := utils.Map(channelsList, toTopic)
	for _, ch := range channelsAsTopics {
		if err := reader.NewSingleChannelConsumer(ch); err != nil {
			return nil, err
		}
	}
	return &reader, nil
}

/*
ReadMessage reads message by message. Returns channel the message belongs to,
the message and an error if any occurred.
*/
func (reader *Reader) ReadMessage(channel string) (models.BrokerData, error) {
	for {
		event := reader.consumers[channel].Poll(pollTimeout)
		switch ev := event.(type) {
		case *kafka.Message:

			channel := *ev.TopicPartition.Topic

			// Decoding Message
			message, errDecode := decode(ev.Value, fromTopic(channel).channel)
			if errDecode != nil {
				return models.BrokerData{}, errDecode
			}

			reader.lastMessage = ev
			channelName := fromTopic(channel).channel

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
func (reader *Reader) CommitMessage(channel string) error {
	_, errCommit := reader.consumers[channel].CommitMessage(reader.lastMessage)
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
	for _, consumer := range reader.consumers {
		err := consumer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (reader *Reader) NewSingleChannelConsumer(channel string) error {
	kafkaEnv := GetEnvironment()
	newConsumer, errKafkaConsumer := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaEnv.KafkaBootstrapServers,
		"group.id":           globalEnv.GetInsprAppID(),
		"auto.offset.reset":  kafkaEnv.KafkaAutoOffsetReset,
		"enable.auto.commit": false,
	})
	if errKafkaConsumer != nil {
		return ierrors.NewError().Message(errKafkaConsumer.Error()).InnerError(errKafkaConsumer).InternalServer().Build()
	}

	newTopic := toTopic(channel)

	if err := newConsumer.Subscribe(newTopic, nil); err != nil {
		return err
	}
	reader.consumers[channel] = newConsumer
	return nil
}
