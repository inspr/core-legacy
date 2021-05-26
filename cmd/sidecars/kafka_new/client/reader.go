package kafkasc

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	globalEnv "github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/sidecars/models"
	"go.uber.org/zap"
)

const pollTimeout = 100

// Reader reads/commit messages from the channels defined in the env
type Reader struct {
	consumers map[string]models.Consumer
}

// NewReader return a new Reader
func NewReader() (models.Reader, error) {
	var reader Reader
	channelsList := globalEnv.GetChannelBoundaryList(globalEnv.GetInputChannels())

	resolvedChList := globalEnv.GetResolvedBoundaryChannelList(globalEnv.GetInputChannels())
	if len(resolvedChList) == 0 {
		return nil, ierrors.NewError().Message("INSPR_INPUT_CHANNELS not specified").InvalidChannel().Build()
	}

	reader.consumers = make(map[string]models.Consumer)

	for idx, ch := range channelsList {
		if err := reader.newSingleChannelConsumer(ch, resolvedChList[idx]); err != nil {
			return nil, err
		}
	}
	return &reader, nil
}

// Consumers returns a Reader's consumers
func (reader *Reader) Consumers() map[string]models.Consumer {
	return reader.consumers
}

/*
ReadMessage reads message by message. Returns channel the message belongs to,
the message and an error if any occurred.
*/
func (reader *Reader) ReadMessage(ctx context.Context, channel string) ([]byte, error) {
	resolved, _ := globalEnv.GetResolvedChannel(channel, globalEnv.GetInputChannels(), "")

	logger.Info("trying to read message from topic",
		zap.String("channel", channel),
		zap.String("resolved channel", resolved),
	)
	consumer := reader.consumers[channel]

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			event := consumer.Poll(pollTimeout)
			switch ev := event.(type) {
			case *kafka.Message:
				topic := *ev.TopicPartition.Topic
				logger.Info("reading message from topic", zap.String("topic", topic))

				return ev.Value, nil

			case kafka.Error:
				if ev.Code() == kafka.ErrAllBrokersDown {
					return nil, ierrors.
						NewError().
						InnerError(ev).
						Message("kafka error = all brokers are down\n%s", ev.Error()).
						InternalServer().
						Build()
				}
				logger.Error("error in reading kafka message", zap.String("error", ev.Error()))
				return nil, ierrors.NewError().
					Message("%v", ev).
					Build()

			default:
				continue
			}

		}
	}
}

// Commit commits the last message read by Reader
func (reader *Reader) Commit(ctx context.Context, channel string) error {
	logger.Info("committing to channel", zap.String("channel", channel))
	doneChan := make(chan error)
	go func() { _, errCommit := reader.consumers[channel].Commit(); doneChan <- errCommit }()
	select {
	case <-ctx.Done():
		<-doneChan
		return ctx.Err()
	case errCommit := <-doneChan:
		if errCommit != nil {
			return ierrors.
				NewError().
				InnerError(errCommit).
				Message("failed to commit last message").
				InternalServer().
				Build()
		}
	}
	return nil
}

// Close closes the reader consumers
func (reader *Reader) Close() error {
	logger.Debug("closing Kafka readers consumers")
	for _, consumer := range reader.consumers {
		err := consumer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//newSingleChannelConsumer creates a consumer for a single Kafka channel on the reader's consumers map.
func (reader *Reader) newSingleChannelConsumer(channel, resolved string) error {
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
