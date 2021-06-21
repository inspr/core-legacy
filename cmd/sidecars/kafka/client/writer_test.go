package kafkasc

import (
	"os"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"inspr.dev/inspr/pkg/environment"
)

func TestNewWriter(t *testing.T) {
	createMockEnv()
	defer deleteMockEnv()
	tests := []struct {
		name    string
		want    *Writer
		wantErr bool
	}{
		{
			name:    "Valid writer creation",
			wantErr: false,
			want:    &Writer{},
		},
		{
			name:    "Invalid writer creation - not mocked (without kafka server up)",
			wantErr: true,
			want:    &Writer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWriter()
			if err != nil {
				t.Error(err)
			}
			defer got.Close()
			if tt.wantErr && (got.Producer().GetFatalError() != nil) {
				t.Errorf("NewWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Producer() == nil {
				t.Errorf("NewWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriter_WriteMessage(t *testing.T) {
	mProd, _ := newMockWriter()
	defer mProd.Close()
	createMockEnv()
	os.Setenv("INSPR_APP_SCOPE", "")
	environment.RefreshEnviromentVariables()
	defer deleteMockEnv()
	type fields struct {
		producer *kafka.Producer
	}
	type args struct {
		channel string
		message []byte
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
				producer: mProd.Producer(),
			},
			args: args{
				channel: "invalid",
				message: []byte("testMessageWriterTest"),
			},
			wantErr: true,
		},
		{
			name: "Valid message writing",
			fields: fields{
				producer: mProd.Producer(),
			},
			args: args{
				channel: "ch1",
				message: []byte("testMessageWriterTest"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &Writer{
				producer: tt.fields.producer,
			}
			if err := writer.WriteMessage(tt.args.channel, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriter_produceMessage(t *testing.T) {
	mProd, _ := newMockWriter()
	defer mProd.Close()
	createMockEnv()
	os.Setenv("INSPR_APP_SCOPE", "")
	environment.RefreshEnviromentVariables()
	defer deleteMockEnv()
	type fields struct {
		producer *kafka.Producer
	}
	type args struct {
		message []byte
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
				producer: mProd.Producer(),
			},
			args: args{
				message: []byte("testProducingMessage"),
				channel: "ch1_resolved",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &Writer{
				producer: tt.fields.producer,
			}
			if err := writer.produceMessage(tt.args.message, tt.args.channel); (err != nil) != tt.wantErr {
				t.Errorf("Writer.produceMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
