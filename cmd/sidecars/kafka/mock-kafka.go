package kafka

import "github.com/confluentinc/confluent-kafka-go/kafka"

//MockConsumer mock
type MockConsumer struct {
	err           error
	errCode       int
	pollMsg       string
	topic         string
	senderChannel string
}

//MockEvent mock
type MockEvent struct {
	message string
}

func (me *MockEvent) String() string {
	return me.message
}

//Poll mock
func (mc *MockConsumer) Poll(timeout int) (event kafka.Event) {

	msg, _ := encode(mc.pollMsg, mc.senderChannel)
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &mc.topic,
		},
		Value: msg,
	}

}

//SubscribeTopics mock
func (mc *MockConsumer) SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) (err error) {
	return nil
}

//CommitMessage mock
func (mc *MockConsumer) CommitMessage(m *kafka.Message) ([]kafka.TopicPartition, error) {
	return nil, nil
}

//Close mock
func (mc *MockConsumer) Close() (err error) {
	return nil
}
