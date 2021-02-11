package kafka

import (
	"os"
	"reflect"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
)

// createMockEnvVars - sets up the env values to be used in the tests functions
// createMockEnvVars - sets up the env values to be used in the tests functions
func createMockReaderEnv() {
	os.Setenv("INSPR_INPUT_CHANNELS", "inp1;inp2;inp3")
	os.Setenv("INSPR_OUTPUT_CHANNELS", "out1;out2;out3")
	os.Setenv("INSPR_UNIX_SOCKET", "/addr/to/socket")
	os.Setenv("INSPR_APP_CTX", "random.app1")
	os.Setenv("INSPR_ENV", "random")
	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "kafka")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "latest")
}

// deleteMockEnvVars - deletes the env values used in the tests functions
func deleteMockReaderEnv() {
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("INSPR_UNIX_SOCKET")
	os.Unsetenv("INSPR_APP_CTX")
	os.Unsetenv("INSPR_ENV")
	os.Unsetenv("KAFKA_BOOTSTRAP_SERVERS")
	os.Unsetenv("KAFKA_AUTO_OFFSET_RESET")
}

func TestNewReader(t *testing.T) {
	createMockReaderEnv()
	defer deleteMockReaderEnv()
	environment.RefreshEnviromentVariables()

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
				deleteMockReaderEnv()
				createMockReaderEnv()
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
	type fields struct {
		consumer    *kafka.Consumer
		lastMessage *kafka.Message
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		want1   interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &Reader{
				consumer:    tt.fields.consumer,
				lastMessage: tt.fields.lastMessage,
			}
			got, got1, err := reader.ReadMessage()
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
		consumer    *kafka.Consumer
		lastMessage *kafka.Message
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &Reader{
				consumer:    tt.fields.consumer,
				lastMessage: tt.fields.lastMessage,
			}
			if err := reader.Commit(); (err != nil) != tt.wantErr {
				t.Errorf("Reader.Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReader_Close(t *testing.T) {
	type fields struct {
		consumer    *kafka.Consumer
		lastMessage *kafka.Message
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &Reader{
				consumer:    tt.fields.consumer,
				lastMessage: tt.fields.lastMessage,
			}
			reader.Close()
		})
	}
}
