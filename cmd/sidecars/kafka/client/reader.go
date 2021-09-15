package kafkasc

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	globalEnv "inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/ierrors"
)

const pollTimeout = 100

// Consumer interface
type Consumer interface {
	Poll(int) kafka.Event
	Commit() ([]kafka.TopicPartition, error)
	Close() (err error)
}

type ReaderMetric struct {
	readKafkaTimeDuration prometheus.Summary
}

// Reader reads/commit messages from the channels defined in the env
type Reader struct {
	consumers map[string]Consumer
	kafkaEnv  *Environment
	metric    map[string]ReaderMetric
}

func (reader *Reader) GetMetric(channel string) ReaderMetric {
	metric, ok := reader.metric[channel]
	if ok {
		return metric
	}
	resolved, _ := globalEnv.GetResolvedChannel(channel, globalEnv.GetInputChannelsData(), globalEnv.GetOutputChannelsData())
	broker := "kafka"
	reader.metric[channel] = ReaderMetric{
		readKafkaTimeDuration: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "kafka",
			Name:      "read_kafka_time_duration",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
			Objectives: map[float64]float64{},
		}),
	}

	return reader.metric[channel]
}

// NewReader return a new Reader
func NewReader() (*Reader, error) {
	logger.Info("creating new kafka reader")
	var reader Reader
	reader.kafkaEnv = GetKafkaEnvironment()
	channelsList := globalEnv.GetChannelBoundaryList(globalEnv.GetInputChannelsData())

	logger.Debug("getting resolved channels list")
	resolvedChList := globalEnv.GetResolvedBoundaryChannelList(globalEnv.GetInputChannelsData())
	if len(resolvedChList) == 0 {
		logger.Error("invalid resolved channel list")
		return nil, ierrors.New(
			"INSPR_INPUT_CHANNELS not specified",
		).InvalidChannel()
	}

	reader.consumers = make(map[string]Consumer)
	reader.metric = make(map[string]ReaderMetric)

	logger.Debug("creating new consumer for each channel")
	for idx, ch := range channelsList {
		if err := reader.newSingleChannelConsumer(ch, resolvedChList[idx]); err != nil {
			logger.Error("unable to create consumer for channel",
				zap.String("channel", ch),
				zap.String("error", err.Error()))

			return nil, err
		}
	}

	logger.Debug("new reader created!")

	newReader := &Reader{
		metric: make(map[string]ReaderMetric),
	}

	return newReader, nil
}

// Consumers returns a Reader's consumers
func (reader *Reader) Consumers() map[string]Consumer {
	return reader.consumers
}

// ReadMessage reads message by message. Returns channel the message belongs to,
// the message and an error if any occurred.
func (reader *Reader) ReadMessage(ctx context.Context, channel string) ([]byte, error) {
	resolved, _ := globalEnv.GetResolvedChannel(channel, globalEnv.GetInputChannelsData(), nil)

	logger.Info("trying to read message from topic",
		zap.String("channel", channel),
		zap.String("resolved channel", resolved),
	)
	consumer := reader.consumers[channel]

	readMsg := time.Now()
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

				elapsed := time.Since(readMsg)
				reader.GetMetric(channel).readKafkaTimeDuration.Observe(elapsed.Seconds())

				return ev.Value, nil

			case kafka.Error:
				if ev.Code() == kafka.ErrAllBrokersDown {
					return nil, ierrors.Wrap(
						ierrors.New(ev).InternalServer(),
						"kafka error = all brokers are down",
					)
				}
				logger.Error("error while reading kafka message", zap.String("error", ev.Error()))
				return nil, ierrors.New("%v", ev)

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
			return ierrors.New(
				"failed to commit last message: %s", errCommit.Error(),
			).InternalServer()
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
	logger.Debug("creating single consumer with configs",
		zap.String("bootstrap", reader.kafkaEnv.KafkaBootstrapServers),
		zap.String("groupid", globalEnv.GetInsprAppID()),
		zap.String("autooffset", reader.kafkaEnv.KafkaAutoOffsetReset))

	newConsumer, errKafkaConsumer := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  reader.kafkaEnv.KafkaBootstrapServers,
		"group.id":           globalEnv.GetInsprAppID(),
		"auto.offset.reset":  reader.kafkaEnv.KafkaAutoOffsetReset,
		"enable.auto.commit": false,
	})
	if errKafkaConsumer != nil {
		return ierrors.New(errKafkaConsumer).InternalServer()
	}

	logger.Debug("subscribing new consumer",
		zap.String("resolved channel", resolved))

	if err := newConsumer.Subscribe(resolved, nil); err != nil {
		return err
	}

	logger.Debug("done subscribing consumer")
	reader.consumers[channel] = newConsumer
	return nil
}
