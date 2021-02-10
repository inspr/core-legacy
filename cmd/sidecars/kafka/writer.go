package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
)

// WriterInterface defines an interface for writing messages
type WriterInterface interface {
	WriteMessage(channel string, message interface{}) error
}

type writer struct {
	producer *kafka.Producer
}

// NewWriter creates a new writer/kafka producer
func NewWriter() (WriterInterface, error) {
	bootstrapServers := GetEnvironment().KafkaBootstrapServers

	kProd, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.server": bootstrapServers,
	})
	if err != nil {
		return nil, kafka.NewError(kafka.ErrInvalidArg, err.Error(), false)
	}

	return &writer{kProd}, nil
}

// WriteMessage receives a message and sends it to the topic defined by the given channel
func (w *writer) WriteMessage(channel string, message interface{}) error {
	if !environment.GetEnvironment().IsInOutputChannel(channel, ";") {
		return kafka.NewError(kafka.ErrInvalidArg, "invalid output channel", false)
	}

	go deliveryReport(w.producer)

	if errProduceMessage := w.produceMessage(message, channel); errProduceMessage != nil {
		return errProduceMessage
	}

	w.producer.Flush(15 * 1000)

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
func (w *writer) produceMessage(message interface{}, channel string) error {
	messageEncoded, errorEncode := encode(message, channel)
	if errorEncode != nil {
		return errorEncode
	}

	topic := toTopic(channel)

	w.producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: messageEncoded,
	}

	return nil
}
