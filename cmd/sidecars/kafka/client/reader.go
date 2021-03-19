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
	CommitMessage(m *kafka.Message) ([]kafka.TopicPartition, error) //deprecated
	Commit() ([]kafka.TopicPartition, error)
	Close() (err error)
}

// Reader reads/commit messages from the channels defined in the env
type Reader struct {
	consumers   map[string]Consumer
	lastMessage *kafka.Message
}

// NewReader return a new Reader
func NewReader() (*Reader, error) {
	var reader Reader
	channelsList := globalEnv.GetInputChannelList(globalEnv.GetInputChannels())
	resolvedChList := globalEnv.GetResolvedInputChannelList(globalEnv.GetInputChannels())
	if len(resolvedChList) == 0 {
		return nil, ierrors.NewError().Message("KAFKA_INPUT_CHANNELS not specified").InvalidChannel().Build()
	}

	reader.consumers = make(map[string]Consumer)

	for idx, ch := range channelsList {
		if err := reader.NewSingleChannelConsumer(ch, resolvedChList[idx]); err != nil {
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
func (reader *Reader) CommitMessage(channel string) error {
	_, errCommit := reader.consumers[channel].Commit()
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

//NewSingleChannelConsumer creates a consumer for a single Kafka channel on the reader's consumers map.
func (reader *Reader) NewSingleChannelConsumer(channel, resolved string) error {
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

	newTopic := messageChannel{channel: resolved}.toTopic()

	if err := newConsumer.Subscribe(newTopic, nil); err != nil {
		return err
	}
	reader.consumers[channel] = newConsumer
	return nil
}
