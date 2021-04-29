package kafkasc

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

//MockConsumer mock
type MockConsumer struct {
	err           bool
	errCode       int
	pollMsg       string
	topic         string
	senderChannel string
	events        chan kafka.Event
}

//MockEvent mock
type MockEvent struct {
	message string
}

func (me *MockEvent) String() string {
	return me.message
}

// Events mock for the events channel
func (mc *MockConsumer) Events() chan kafka.Event {
	return mc.events
}

// CreateEvent creates an event on the mock channel of thre consumer
func (mc *MockConsumer) CreateEvent(ev kafka.Event) {
	mc.events <- ev
}

// CreateMessage creates a message on the mock consumer from the field defined
// in the structure
func (mc *MockConsumer) CreateMessage() {

	if mc.err {
		mc.events <- kafka.NewError(kafka.ErrAllBrokersDown, "", false)
	}

	ch := kafkaTopic(mc.senderChannel)
	msg, _ := ch.encode(mc.pollMsg)
	mc.events <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &mc.topic,
		},
		Value: msg,
	}
}

//Commit mock
func (mc *MockConsumer) Commit() ([]kafka.TopicPartition, error) {
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
	os.Setenv("ch1_resolved_SCHEMA", `{"type":"string"}`)
	os.Setenv("ch2_resolved_SCHEMA", "hellotest")
	os.Setenv("ch1_RESOLVED", `ch1_resolved`)
	os.Setenv("ch2_RESOLVED", "ch2_resolved")
	os.Setenv("INSPR_APP_ID", "testappid1")
	os.Setenv("INSPR_SIDECAR_IMAGE", "random-sidecar-image")
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
	os.Unsetenv("INSPR_APP_ID")
	os.Unsetenv("INSPR_SIDECAR_IMAGE")

	os.Unsetenv("ch1_RESOLVED")
	os.Unsetenv("ch2_RESOLVED")
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
