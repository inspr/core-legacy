package kafka

import (
	"os"
	"reflect"
	"testing"

	"github.com/linkedin/goavro"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
)

var mockSchema = `{"type":"string"}`

func mockNewCodec() *goavro.Codec {
	codec, _ := goavro.NewCodec(mockSchema)
	return codec
}

func returnEncodedMessage(msg string) []byte {
	bMsg, _ := mockNewCodec().BinaryFromNative(nil, msg)
	return bMsg
}

func Test_getCodec(t *testing.T) {
	type args struct {
		schema string
	}
	tests := []struct {
		name    string
		args    args
		want    *goavro.Codec
		wantErr bool
	}{
		{
			name: "Valid schema to generate new codec",
			args: args{
				schema: mockSchema,
			},
			want:    mockNewCodec(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getMockApp()
			got, err := getCodec(tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCodec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Schema() != tt.want.Schema() {
				t.Errorf("getCodec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSchema(t *testing.T) {
	type args struct {
		channel string
		context string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Invalid channel",
			args: args{
				channel: "invalid",
				context: "",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Valid channel with schema",
			args: args{
				channel: "ch2",
				context: "",
			},
			want:    "hellotest",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getMockApp()
			got, err := getSchema(tt.args.channel, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decode(t *testing.T) {
	createMockReaderEnv()
	os.Setenv("INSPR_APP_CTX", "")
	environment.RefreshEnviromentVariables()
	defer deleteMockReaderEnv()
	type args struct {
		messageEncoded []byte
		channel        string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Invalid channel",
			args: args{
				channel:        "invalid",
				messageEncoded: []byte{},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Invalid schema",
			args: args{
				channel:        "ch2",
				messageEncoded: []byte{104, 101, 108, 108, 111, 116, 101, 115, 116},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Valid schema",
			args: args{
				channel:        "ch1",
				messageEncoded: returnEncodedMessage("testSchemaString"),
			},
			wantErr: false,
			want:    "testSchemaString",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getMockApp()
			got, err := decode(tt.args.messageEncoded, tt.args.channel)
			if (err != nil) != tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encode(t *testing.T) {
	createMockReaderEnv()
	os.Setenv("INSPR_APP_CTX", "")
	environment.RefreshEnviromentVariables()
	defer deleteMockReaderEnv()
	type args struct {
		message interface{}
		channel string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Invalid channel",
			args: args{
				channel: "invalid",
				message: []byte{},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Invalid schema",
			args: args{
				channel: "ch2",
				message: []byte{104, 101, 108, 108, 111, 116, 101, 115, 116},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Valid encoding",
			args: args{
				channel: "ch1",
				message: "testMessageEncodingString",
			},
			wantErr: false,
			want:    returnEncodedMessage("testMessageEncodingString"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getMockApp()
			got, err := encode(tt.args.message, tt.args.channel)
			if (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
