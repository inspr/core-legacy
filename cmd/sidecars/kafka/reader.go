package kafka

import (
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	globalEnv "gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// Reader reads/commit messages from the channels defined in the env
type Reader struct {
	consumer    *kafka.Consumer
	lastMessage *kafka.Message
}

// NewReader return a new Reader
func NewReader() (*Reader, error) {
	kafkaEnv := GetEnvironment()
	globalEnv := globalEnv.GetEnvironment()

	var reader Reader

	newConsumer, errKafkaConsumer := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaEnv.KafkaBootstrapServers,
		"group.id":           globalEnv.InsprAppContext,
		"auto.offset.reset":  kafkaEnv.KafkaAutoOffsetReset,
		"enable.auto.commit": false,
	})
	if errKafkaConsumer != nil {
		return nil, ierrors.NewError().Message("failed to create a new kafka consumer").InnerError(errKafkaConsumer).InternalServer().Build()
	}
	reader.consumer = newConsumer

	channelsList := globalEnv.InputChannels
	if len(channelsList) == 0 {
		return nil, ierrors.NewError().Message("KAFKA_INPUT_CHANNELS not specified").InternalServer().Build()
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
func (reader *Reader) ReadMessage() (string, interface{}, error) {
	for {
		event := reader.consumer.Poll(100)
		switch ev := event.(type) {
		case *kafka.Message:

			channel := *ev.TopicPartition.Topic

			// Decoding Message
			message, errDecode := decode(ev.Value, fromTopic(channel).channel)
			if errDecode != nil {
				return "", nil, errDecode
			}

			reader.lastMessage = ev
			channelName := fromTopic(channel).channel
			return channelName, message, nil

		case kafka.Error:
			if ev.Code() == kafka.ErrAllBrokersDown {
				return "", nil, ierrors.
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
func (reader *Reader) Commit() error {
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
func (reader *Reader) Close() {
	reader.consumer.Close()
}
