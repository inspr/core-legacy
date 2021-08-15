package kafkasc

import (
	"os"
	"reflect"
	"testing"
)

func mockKafkaEnvironment() *Environment {
	return &Environment{
		KafkaBootstrapServers: "localhost",
		KafkaAutoOffsetReset:  "101019",
	}
}

func TestGetEnvironment(t *testing.T) {
	os.Setenv("INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS", "localhost")
	os.Setenv("INSPR_SIDECAR_KAFKA_AUTO_OFFSET_RESET", "101019")
	defer deleteMockEnv()
	tests := []struct {
		name string
		want *Environment
	}{
		{
			name: "Get all environment variables",
			want: mockKafkaEnvironment(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKafkaEnvironment(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnv(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get bootstrap servers environment variable",
			args: args{
				name: "INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS",
			},
			want: "localhost",
		},
		{
			name: "Get auto offset reset environment variable",
			args: args{
				name: "INSPR_SIDECAR_KAFKA_AUTO_OFFSET_RESET",
			},
			want: "101019",
		},
		{
			name: "Invalid - Get invalid environment variable",
			args: args{
				name: "KAFKA_INVALID_ENV_VAR",
			},
			want: "[ENV VAR] KAFKA_INVALID_ENV_VAR not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer func() {
				recover()
			}()

			if got := getEnv(tt.args.name); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefreshEnviromentVariables(t *testing.T) {
	createMockEnv()
	os.Setenv("INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS", "one")
	os.Setenv("INSPR_SIDECAR_KAFKA_AUTO_OFFSET_RESET", "two")
	defer deleteMockEnv()
	tests := []struct {
		name    string
		refresh bool
		want    *Environment
	}{
		{
			name:    "Changed and refreshed environment variables",
			refresh: true,
			want: &Environment{
				KafkaBootstrapServers: "one",
				KafkaAutoOffsetReset:  "two",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := RefreshEnviromentVariables(); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("GetEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
