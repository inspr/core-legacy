package kafkasc

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/ierrors"
)

const flushTimeout = 1000

type writerMetrics struct {
	resolveChannelDuration prometheus.Summary
	produceMessageDuration prometheus.Summary
	flushDuration          prometheus.Summary
}

// Writer defines an interface for writing messages
type Writer struct {
	producer *kafka.Producer
	metrics  map[string]writerMetrics
}

// NewWriter creates a new writer/kafka producer
func NewWriter() (*Writer, error) {
	var kProd *kafka.Producer
	var err error

	bootstrapServers := GetKafkaEnvironment().KafkaBootstrapServers
	kProd, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})

	go func(events <-chan kafka.Event) {
		for {
			switch ev := (<-events).(type) {
			case *kafka.Error:
				logger.Warn("kafka has created an error event", zap.Error(ev))
			case *kafka.LogEvent:
				logger.Debug(ev.Message, zap.String("sub-section", ev.Tag))
			default:
			}
		}
	}(kProd.Events())

	if err != nil {
		return nil, ierrors.New(err)
	}

	newWriter := &Writer{
		producer: kProd,
		metrics:  make(map[string]writerMetrics),
	}

	return newWriter, nil
}

func (writer *Writer) getProducer() *kafka.Producer {
	return writer.producer
}

// WriteMessage receives a message and sends it to the topic defined by the given channel
func (writer *Writer) WriteMessage(channel string, message []byte) error {
	outputChan := environment.GetOutputChannelsData()

	startResolveChannel := time.Now()

	resolvedCh, err := environment.GetResolvedChannel(channel, nil, outputChan)
	if err != nil {
		return err
	}

	elapsedResolveChannel := time.Since(startResolveChannel)
	writer.GetMetric(channel).resolveChannelDuration.Observe(elapsedResolveChannel.Seconds())

	startProduce := time.Now()

	logger.Info("trying to write message in topic",
		zap.String("channel", channel),
		zap.String("resolved channel", resolvedCh))

	if errProduceMessage := writer.produceMessage(message, resolvedCh); errProduceMessage != nil {
		logger.Error("error while producing message",
			zap.Any("error", errProduceMessage))
		return errProduceMessage
	}

	elapsedProduce := time.Since(startProduce)
	writer.GetMetric(channel).produceMessageDuration.Observe(elapsedProduce.Seconds())

	startFlush := time.Now()

	logger.Debug("flushing the producer")
	writer.producer.Flush(flushTimeout)
	logger.Debug("flushed")

	elapsedFlush := time.Since(startFlush)
	writer.GetMetric(channel).flushDuration.Observe(elapsedFlush.Seconds())

	return nil
}

// creates a Kafka message and sends it through the ProduceChannel
func (writer *Writer) produceMessage(message []byte, resolvedChannel string) error {

	logger.Debug("writing message into Kafka Topic",
		zap.String("topic", resolvedChannel))

	return writer.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &resolvedChannel,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}, nil)

}

// Close closes the kafka producer
func (writer *Writer) Close() {
	logger.Debug("closing Kafka producer")
	writer.producer.Close()
}

func (writer *Writer) GetMetric(channel string) writerMetrics {
	metric, ok := writer.metrics[channel]
	if ok {
		return metric
	}
	resolved, _ := environment.GetResolvedChannel(channel, environment.GetInputChannelsData(), environment.GetOutputChannelsData())
	broker := "kafka"
	writer.metrics[channel] = writerMetrics{
		resolveChannelDuration: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "kafka_sidecar_writer",
			Name:      "resolve_channel_duration",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
			Objectives: map[float64]float64{},
		}),
		produceMessageDuration: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "kafka_sidecar_writer",
			Name:      "produce_message_duration",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
			Objectives: map[float64]float64{},
		}),
		flushDuration: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "kafka_sidecar_writer",
			Name:      "flush_duration",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
			Objectives: map[float64]float64{},
		}),
	}

	return writer.metrics[channel]

}
