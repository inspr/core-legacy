package kafkasc

import (
	"os"
	"reflect"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"go.uber.org/zap"
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
				if !(reader.consumers != nil && len(reader.consumers) > 0) {
					t.Errorf("check function error = Reader not created successfully")
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
		consumers   map[string]Consumer
		lastMessage *kafka.Message
	}
	tests := []struct {
		name          string
		fields        fields
		before        func()
		uniqueChannel string
		want          string
		want1         interface{}
		wantErr       bool
	}{
		{
			name: "It should read a message",
			fields: fields{
				consumers: map[string]Consumer{
					"ch1_resolved": &MockConsumer{
						err:           false,
						pollMsg:       "Hello World!",
						topic:         "ch1_resolved",
						errCode:       0,
						senderChannel: "ch1_resolved",
					},
				},
				lastMessage: nil,
			},
			wantErr:       false,
			uniqueChannel: "ch1_resolved",
			want:          "ch1_resolved",
			want1:         "Hello World!",
		},
		{
			name: "It should return a message poll error",
			fields: fields{
				consumers: map[string]Consumer{
					"ch1_resolved": &MockConsumer{
						err:           true,
						pollMsg:       "Hello World!",
						topic:         "ch1_resolved",
						errCode:       0,
						senderChannel: "ch1_resolved",
					},
				},
				lastMessage: nil,
			},
			uniqueChannel: "ch1_resolved",
			wantErr:       true,
		},
		{
			name: "It should return a decode error (sender channel invalid)",
			fields: fields{
				consumers: map[string]Consumer{
					"ch1_resolved": &MockConsumer{
						err:           false,
						pollMsg:       "Hello World!",
						topic:         "ch1_resolved",
						errCode:       0,
						senderChannel: "ch2",
					},
				},
				lastMessage: nil,
			},
			uniqueChannel: "ch1_resolved",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, _ := zap.NewDevelopment()
			reader := &Reader{
				consumers: tt.fields.consumers,
				logger:    logger,
			}

			bData, err := reader.ReadMessage(tt.uniqueChannel)
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
		consumers   map[string]Consumer
		lastMessage *kafka.Message
	}
	tests := []struct {
		name          string
		fields        fields
		uniqueChannel string
		wantErr       bool
	}{
		{
			name: "It should not return a error since the message was committed",
			fields: fields{
				consumers: map[string]Consumer{
					"ch1": &MockConsumer{
						err: false,
					},
				},
				lastMessage: nil,
			},
			uniqueChannel: "ch1",
			wantErr:       false,
		},
		{
			name: "It should return a error since the message was not committed",
			fields: fields{
				consumers: map[string]Consumer{
					"ch1": &MockConsumer{
						err: true,
					},
				},
				lastMessage: nil,
			},
			uniqueChannel: "ch1",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &Reader{
				consumers: tt.fields.consumers,
			}
			if err := reader.Commit(tt.uniqueChannel); (err != nil) != tt.wantErr {
				t.Errorf("Reader.Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReader_Close(t *testing.T) {
	type fields struct {
		consumers   map[string]Consumer
		lastMessage *kafka.Message
	}
	tests := []struct {
		name          string
		fields        fields
		uniqueChannel string
		wantErr       bool
	}{
		{
			name: "Close the consumer",
			fields: fields{
				consumers: map[string]Consumer{
					"ch1": &MockConsumer{
						err: false,
					},
				},
				lastMessage: nil,
			},
			wantErr: false,
		},
		{
			name: "Error when trying to close the consumer",
			fields: fields{
				consumers: map[string]Consumer{
					"ch1": &MockConsumer{
						err: true,
					},
				},
				lastMessage: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &Reader{
				consumers: tt.fields.consumers,
			}
			if err := reader.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Reader.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
