package kafkasc

import (
	"os"
	"reflect"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
)

func TestNewReader(t *testing.T) {
	createMockEnv()
	defer deleteMockEnv()
	environment.RefreshEnviromentVariables()
	RefreshEnviromentVariables()

	tests := []struct {
		name          string
		want          *Reader
		wantErr       bool
		checkFunction func(t *testing.T, reader *Reader)
		before        func()
	}{
		{
			name:    "It should return a new Reader",
			wantErr: false,
			checkFunction: func(t *testing.T, reader *Reader) {
				if !(reader.consumer != nil && reader.lastMessage == nil) {
					t.Errorf("check function error = Reader not created sucesfully")
				}
			},
		},
		{
			name:    "Input channel list is empty - it should return a error",
			wantErr: true,
			before: func() {
				environment.RefreshEnviromentVariables()
				deleteMockEnv()
				createMockEnv()
				os.Setenv("INSPR_INPUT_CHANNELS", "")
				environment.RefreshEnviromentVariables()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}
			got, err := NewReader()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunction != nil {
				tt.checkFunction(t, got)
			}
		})
	}
}

func TestReader_ReadMessage(t *testing.T) {
	createMockEnv()
	defer deleteMockEnv()
	environment.RefreshEnviromentVariables()
	RefreshEnviromentVariables()

	type fields struct {
		consumer    Consumer
		lastMessage *kafka.Message
	}
	tests := []struct {
		name    string
		fields  fields
		before  func()
		want    string
		want1   interface{}
		wantErr bool
	}{
		{
			name: "It should read a message",
			fields: fields{
				consumer: &MockConsumer{
					err:           false,
					pollMsg:       "Hello World!",
					topic:         toTopic("ch1"),
					errCode:       0,
					senderChannel: "ch1",
				},
				lastMessage: nil,
			},
			wantErr: false,
			want:    "ch1",
			want1:   "Hello World!",
		},
		{
			name: "It should return a message poll error",
			fields: fields{
				consumer: &MockConsumer{
					err:           true,
					pollMsg:       "Hello World!",
					topic:         toTopic("ch1"),
					errCode:       0,
					senderChannel: "ch1",
				},
				lastMessage: nil,
			},
			wantErr: true,
		},
		{
			name: "It should return a decode error (sender channel invalid)",
			fields: fields{
				consumer: &MockConsumer{
					err:           false,
					pollMsg:       "Hello World!",
					topic:         toTopic("ch1"),
					errCode:       0,
					senderChannel: "ch2",
				},
				lastMessage: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &Reader{
				consumer:    tt.fields.consumer,
				lastMessage: tt.fields.lastMessage,
			}

			bData, err := reader.ReadMessage()
			got := bData.Channel
			got1 := bData.Message.Data

			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Reader.ReadMessage() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Reader.ReadMessage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestReader_Commit(t *testing.T) {
	type fields struct {
		consumer    Consumer
		lastMessage *kafka.Message
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "It should not return a error since the message was commited",
			fields: fields{
				consumer: &MockConsumer{
					err: false,
				},
				lastMessage: nil,
			},
			wantErr: false,
		},
		{
			name: "It should return a error since the message was not commited",
			fields: fields{
				consumer: &MockConsumer{
					err: true,
				},
				lastMessage: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &Reader{
				consumer:    tt.fields.consumer,
				lastMessage: tt.fields.lastMessage,
			}
			if err := reader.CommitMessage(); (err != nil) != tt.wantErr {
				t.Errorf("Reader.Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReader_Close(t *testing.T) {
	type fields struct {
		consumer    Consumer
		lastMessage *kafka.Message
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Close the consumer",
			fields: fields{
				consumer:    &MockConsumer{},
				lastMessage: nil,
			},
			wantErr: false,
		},
		{
			name: "Error when trying to close the consumer",
			fields: fields{
				consumer: &MockConsumer{
					err: true,
				},
				lastMessage: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &Reader{
				consumer:    tt.fields.consumer,
				lastMessage: tt.fields.lastMessage,
			}
			if err := reader.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Reader.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
