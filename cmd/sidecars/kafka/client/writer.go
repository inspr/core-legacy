package kafkasc

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/ierrors"
)

const flushTimeout = 1000

// Writer defines an interface for writing messages
type Writer struct {
	producer *kafka.Producer
}

// NewWriter creates a new writer/kafka producer
func NewWriter() (*Writer, error) {
	var kProd *kafka.Producer
	var err error

	bootstrapServers := GetKafkaEnvironment().KafkaBootstrapServers
	kProd, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})

	if err != nil {
		return nil, ierrors.NewError().Message(err.Error()).Build()
	}

	return &Writer{kProd}, nil
}

// Producer returns a Writer's producer
func (writer *Writer) Producer() *kafka.Producer {
	return writer.producer
}

// WriteMessage receives a message and sends it to the topic defined by the given channel
func (writer *Writer) WriteMessage(channel string, message []byte) error {
	outputChan := environment.GetOutputChannelsData()

	resolvedCh, err := environment.GetResolvedChannel(channel, nil, outputChan)
	if err != nil {
		return err
	}

	logger.Info("trying to write message in topic",
		zap.String("channel", channel),
		zap.String("resolved channel", resolvedCh))

	if errProduceMessage := writer.produceMessage(message, resolvedCh); errProduceMessage != nil {
		logger.Error("error while producing message",
			zap.Any("error", errProduceMessage))
		return errProduceMessage
	}

	logger.Debug("flushing the producer")
	writer.producer.Flush(flushTimeout)
	logger.Debug("flushed")
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
