package kafkasc

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

//MockConsumer mock
type MockConsumer struct {
	err           bool
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

	if mc.err {
		return kafka.NewError(kafka.ErrAllBrokersDown, "", false)
	}

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
	if mc.err {
		return nil, kafka.NewError(kafka.ErrApplication, "", false)
	}
	return nil, nil
}

//Close mock
func (mc *MockConsumer) Close() (err error) {
	if mc.err {
		return kafka.NewError(kafka.ErrApplication, "", false)
	}
	return nil
}

// createMockEnvVars - sets up the env values to be used in the tests functions
// createMockEnvVars - sets up the env values to be used in the tests functions
func createMockEnv() {
	os.Setenv("INSPR_INPUT_CHANNELS", "ch1;ch2")
	os.Setenv("INSPR_OUTPUT_CHANNELS", "ch1;ch2")
	os.Setenv("INSPR_UNIX_SOCKET", "/addr/to/socket")
	os.Setenv("INSPR_APP_CTX", "")
	os.Setenv("INSPR_ENV", "random")
	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "kafka")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "latest")
}

// deleteMockEnvVars - deletes the env values used in the tests functions
func deleteMockEnv() {
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("INSPR_UNIX_SOCKET")
	os.Unsetenv("INSPR_APP_CTX")
	os.Unsetenv("INSPR_ENV")
	os.Unsetenv("KAFKA_BOOTSTRAP_SERVERS")
	os.Unsetenv("KAFKA_AUTO_OFFSET_RESET")
}

// mockMessageSender sends two messages to a kafka producer
func mockMessageSender(writer *kafka.Producer, topic *string) {
	// Valid message
	writer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte("msgTest"),
	}
	// Invalid message
	writer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{},
		Value:          []byte("msgTest"),
	}
}

// getMockApp returns a mocked dApp
func getMockApp() tree.MemoryManager {
	ctype := &meta.ChannelType{
		Meta: meta.Metadata{
			Name:        "ct1",
			Reference:   "root.ct1",
			Annotations: map[string]string{},
			Parent:      "root",
			SHA256:      "",
		},
		Schema: []byte{123, 34, 116, 121, 112, 101, 34, 58, 34, 115, 116, 114, 105, 110, 103, 34, 125},
	}
	ctype1 := &meta.ChannelType{
		Meta: meta.Metadata{
			Name:        "ct2",
			Reference:   "root.ct2",
			Annotations: map[string]string{},
			Parent:      "root",
			SHA256:      "",
		},
		Schema: []byte{104, 101, 108, 108, 111, 116, 101, 115, 116},
	}
	chann := &meta.Channel{
		Meta: meta.Metadata{
			Name:        "ch1",
			Reference:   "root.ch1",
			Annotations: map[string]string{},
			Parent:      "root",
			SHA256:      "",
		},
		Spec: meta.ChannelSpec{
			Type: "ct1",
		},
	}
	chann1 := &meta.Channel{
		Meta: meta.Metadata{
			Name:        "ch2",
			Reference:   "root.ch2",
			Annotations: map[string]string{},
			Parent:      "root",
			SHA256:      "",
		},
		Spec: meta.ChannelSpec{
			Type: "ct2",
		},
	}
	tree.GetTreeMemory()
	tree.GetTreeMemory().ChannelTypes().CreateChannelType(ctype, "")
	tree.GetTreeMemory().ChannelTypes().CreateChannelType(ctype1, "")
	tree.GetTreeMemory().Channels().CreateChannel("", chann)
	tree.GetTreeMemory().Channels().CreateChannel("", chann1)
	return tree.MemoryManager{}
}
