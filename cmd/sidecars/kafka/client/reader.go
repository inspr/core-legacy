package kafkasc

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	globalEnv "gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
	"go.uber.org/zap"
)

const pollTimeout = 100

// Consumer interface
type Consumer interface {
	Poll(timeout int) (event kafka.Event)
	Commit() ([]kafka.TopicPartition, error)
	Close() (err error)
}

// Reader reads/commit messages from the channels defined in the env
type Reader struct {
	consumers map[string]Consumer
	logger    *zap.Logger
}

// NewReader return a new Reader
func NewReader() (*Reader, error) {
	var reader Reader
	channelsList := globalEnv.GetChannelBoundaryList(globalEnv.GetInputChannels())
	resolvedChList := globalEnv.GetResolvedBoundaryChannelList(globalEnv.GetInputChannels())

	logger, _ := zap.NewDevelopment(zap.Fields(zap.String("section", "kafka sidecar")))
	if len(resolvedChList) == 0 {
		return nil, ierrors.NewError().Message("KAFKA_INPUT_CHANNELS not specified").InvalidChannel().Build()
	}
	reader.logger = logger
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
	resolved, _ := globalEnv.GetResolvedChannel(channel, globalEnv.GetInputChannels(), "")
	reader.logger.Info("trying to read message from topic",
		zap.String("channel", channel),
		zap.String("resolved channel", resolved),
	)
	consumer := reader.consumers[channel]

	for {
		event := consumer.Poll(pollTimeout)
		switch ev := event.(type) {
		case *kafka.Message:
			topic := *ev.TopicPartition.Topic
			reader.logger.Info("reading message from topic", zap.String("topic", topic))
			// converts string to defined type for helper functions
			return kafkaTopic(topic).readMessage(ev.Value)

		case kafka.Error:
			if ev.Code() == kafka.ErrAllBrokersDown {
				return models.BrokerData{}, ierrors.
					NewError().
					InnerError(ev).
					Message("kafka error = all brokers are down\n%s", ev.Error()).
					InternalServer().
					Build()
			}
			reader.logger.Error("error in reading kafka message", zap.String("error", ev.Error()))
			return models.BrokerData{}, ierrors.NewError().
				Message("%v", ev).
				Build()

		default:
			continue
		}
	}
}

// Commit commits the last message read by Reader
func (reader *Reader) Commit(channel string) error {
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

	if err := newConsumer.Subscribe(resolved, nil); err != nil {
		return err
	}
	reader.consumers[channel] = newConsumer
	return nil
}
