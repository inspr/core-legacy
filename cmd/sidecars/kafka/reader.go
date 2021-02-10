package kafka

import (
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	globalEnv "gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

type readerObj struct {
	consumer    *kafka.Consumer
	lastMessage *kafka.Message
}

// Reader reads messages from the channels defined in the env
type Reader interface {
	Commit() error
	ReadMessage() (*string, interface{}, error)
	Close()
}

// NewReader return a new reader
func NewReader() (Reader, error) {
	kafkaEnv := GetEnvironment()
	globalEnv := globalEnv.GetEnvironment()

	var reader readerObj

	newConsumer, errKafkaConsumer := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaEnv.KafkaBootstrapServers,
		"group.id":           globalEnv.InsprAppContext,
		"auto.offset.reset":  kafkaEnv.KafkaAutoOffsetReset,
		"enable.auto.commit": false,
	})
	if errKafkaConsumer != nil {
		return nil, ierrors.NewError().Message("[NEW READER] " + errKafkaConsumer.Error()).InternalServer().Build()
	}
	reader.consumer = newConsumer

	channelsList := globalEnv.InputChannels
	if len(channelsList) == 0 {
		return nil, ierrors.NewError().Message("[ENV VAR] KAFKA_INPUT_CHANNELS not specified").InternalServer().Build()
	}

	channelList := strings.Split(channelsList, ";")
	channelList = channelList[:len(channelList)-1]
	channelsToConsume := func() []string {
		ret := []string{}
		for _, s := range channelList {
			topic := toTopic(s)
			ret = append(ret, topic)
		}
		return ret
	}()

	reader.consumer.SubscribeTopics(channelsToConsume, nil)
	return &reader, nil
}

/*
ReadMessage reads message by message. Returns channel the message belongs to,
the message and an error if any occured.
*/
func (reader *readerObj) ReadMessage() (*string, interface{}, error) {
	for {
		event := reader.consumer.Poll(100)
		switch ev := event.(type) {
		case *kafka.Message:

			channel := *ev.TopicPartition.Topic

			// Decoding Message
			message, errDecode := decode(ev.Value, fromTopic(channel).channel)
			if errDecode != nil {
				return nil, nil, errDecode
			}

			reader.lastMessage = ev
			channelName := fromTopic(channel).channel
			return &channelName, message, nil

		case kafka.Error:
			if ev.Code() == kafka.ErrAllBrokersDown {
				return nil, nil, ierrors.
					NewError().
					Message("[FATAL_ERROR]\n===== All brokers are down! =====\n" + ev.Error()).
					InternalServer().
					Build()
			}

		default:
			continue
		}
	}
}

// Commit commits the last message read by Reader
func (reader *readerObj) Commit() error {
	_, errCommit := reader.consumer.CommitMessage(reader.lastMessage)
	if errCommit != nil {
		return ierrors.
			NewError().
			Message("[READER_COMMIT] " + errCommit.Error()).
			InternalServer().
			Build()
	}
	return nil
}

// Close close the reader consumer
func (reader *readerObj) Close() {
	reader.consumer.Close()
}
