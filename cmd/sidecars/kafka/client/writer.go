package kafkasc

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"inspr.dev/inspr/pkg/environment"
	"go.uber.org/zap"
)

const flushTimeout = 15 * 1000

// Writer defines an interface for writing messages
type Writer struct {
	producer *kafka.Producer
}

// NewWriter creates a new writer/kafka producer
func NewWriter(mock bool) (*Writer, error) {
	var kProd *kafka.Producer
	var err error
	if mock {
		kProd, _ = kafka.NewProducer(&kafka.ConfigMap{
			"test.mock.num.brokers": 3,
		})
	} else {
		bootstrapServers := GetEnvironment().KafkaBootstrapServers
		kProd, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": bootstrapServers,
		})
		if err != nil {
			return nil, kafka.NewError(kafka.ErrInvalidArg, err.Error(), false)
		}
	}

	return &Writer{kProd}, nil
}

// WriteMessage receives a message and sends it to the topic defined by the given channel
func (writer *Writer) WriteMessage(channel string, message interface{}) error {
	outputChan := environment.GetOutputChannels()

	resolvedCh, _ := environment.GetResolvedChannel(channel, "", outputChan)

	logger.Info("trying to write message in topic",
		zap.String("channel", channel),
		zap.String("resolved channel", resolvedCh))
	go deliveryReport(writer.producer)

	if errProduceMessage := writer.produceMessage(message, kafkaTopic(resolvedCh)); errProduceMessage != nil {
		logger.Error("error while producing message",
			zap.Any("error", errProduceMessage))
		return errProduceMessage
	}

	writer.producer.Flush(flushTimeout)

	return nil
}

// Logs the ProduceChannel events for successful and failed messages sent
func deliveryReport(producer *kafka.Producer) {
	for event := range producer.Events() {
		switch ev := event.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				log.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
			return
		default:
			log.Println(ev)
		}
	}
}

// creates a Kafka message and sends it through the ProduceChannel
func (writer *Writer) produceMessage(message interface{}, resolvedChannel kafkaTopic) error {
	messageEncoded, errorEncode := resolvedChannel.encode(message)
	if errorEncode != nil {
		return errorEncode
	}

	logger.Debug("writing message into Kafka Topic",
		zap.String("topic", string(resolvedChannel)))

	writer.producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     (*string)(&resolvedChannel),
			Partition: kafka.PartitionAny,
		},
		Value: messageEncoded,
	}

	return nil
}

// Close closes the kafka producer
func (writer *Writer) Close() {
	logger.Debug("closing Kafka producer")
	writer.producer.Close()
}
