package kafka

import (
	"os"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
)

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
func TestNewWriter(t *testing.T) {
	createMockEnvVars()
	defer deleteMockEnvVars()
	type args struct {
		mock bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Writer
		wantErr bool
	}{
		{
			name: "Valid writer creation",
			args: args{
				mock: true,
			},
			wantErr: false,
			want:    &Writer{},
		},
		{
			name: "Invalid writer creation - not mocked (without kafka server up)",
			args: args{
				mock: false,
			},
			wantErr: true,
			want:    &Writer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWriter(tt.args.mock)
			if tt.wantErr && (got.producer.GetFatalError() != nil) {
				t.Errorf("NewWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.producer == nil {
				t.Errorf("NewWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriter_WriteMessage(t *testing.T) {
	mProd, _ := NewWriter(true)
	createMockEnvVars()
	os.Setenv("INSPR_APP_CTX", "")
	environment.RefreshEnviromentVariables()
	defer deleteMockEnvVars()
	type fields struct {
		producer *kafka.Producer
	}
	type args struct {
		channel string
		message interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Invalid channel",
			fields: fields{
				producer: mProd.producer,
			},
			args: args{
				channel: "invalid",
				message: "testMessageWriterTest",
			},
			wantErr: true,
		},
		{
			name: "Valid message writing",
			fields: fields{
				producer: mProd.producer,
			},
			args: args{
				channel: "ch1",
				message: "testMessageWriterTest",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getMockApp()
			writer := &Writer{
				producer: tt.fields.producer,
			}
			if err := writer.WriteMessage(tt.args.channel, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_deliveryReport(t *testing.T) {
	mProd, _ := NewWriter(true)
	type args struct {
		producer *kafka.Producer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Read two messages from de producer",
			args: args{
				producer: mProd.producer,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var strPointer = new(string)
			*strPointer = "testTopic"
			mockMessageSender(tt.args.producer, strPointer)
			for i := 0; i < 2; i++ {
				deliveryReport(tt.args.producer)
			}
		})
	}
}

func TestWriter_produceMessage(t *testing.T) {
	mProd, _ := NewWriter(true)
	createMockEnvVars()
	os.Setenv("INSPR_APP_CTX", "")
	environment.RefreshEnviromentVariables()
	defer deleteMockEnvVars()
	type fields struct {
		producer *kafka.Producer
	}
	type args struct {
		message interface{}
		channel string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid production of given message",
			fields: fields{
				producer: mProd.producer,
			},
			args: args{
				message: "testProducingMessage",
				channel: "ch1",
			},
			wantErr: false,
		},
		{
			name: "Invalid production - encode error",
			fields: fields{
				producer: mProd.producer,
			},
			args: args{
				message: "testProducingMessage",
				channel: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getMockApp()
			writer := &Writer{
				producer: tt.fields.producer,
			}
			if err := writer.produceMessage(tt.args.message, tt.args.channel); (err != nil) != tt.wantErr {
				t.Errorf("Writer.produceMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
