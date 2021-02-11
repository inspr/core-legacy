package kafka

import (
	"os"
	"reflect"
	"testing"
)

// createMockEnvVars - sets up the env values to be used in the tests functions
// createMockEnvVars - sets up the env values to be used in the tests functions
func createMockEnvVars() {
	os.Setenv("INSPR_INPUT_CHANNELS", "inp1;inp2;inp3;")
	os.Setenv("INSPR_OUTPUT_CHANNELS", "out1;out2;out3;ch1;")
	os.Setenv("INSPR_UNIX_SOCKET", "/addr/to/socket")
	os.Setenv("INSPR_APP_CTX", "random.app1")
	os.Setenv("INSPR_ENV", "random")
	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "localhost")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "101019")
}

// deleteMockEnvVars - deletes the env values used in the tests functions
func deleteMockEnvVars() {
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("INSPR_UNIX_SOCKET")
	os.Unsetenv("INSPR_APP_CTX")
	os.Unsetenv("INSPR_ENV")
	os.Unsetenv("KAFKA_BOOTSTRAP_SERVERS")
	os.Unsetenv("KAFKA_AUTO_OFFSET_RESET")
}

func mockKafkaEnvironment() *Environment {
	return &Environment{
		KafkaBootstrapServers: "localhost",
		KafkaAutoOffsetReset:  "101019",
	}
}

func TestGetEnvironment(t *testing.T) {
	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "localhost")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "101019")
	defer deleteMockEnvVars()
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
			if got := GetEnvironment(); !reflect.DeepEqual(got, tt.want) {
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
			name: "Get bootstrap servers enviroment variable",
			args: args{
				name: "KAFKA_BOOTSTRAP_SERVERS",
			},
			want: "localhost",
		},
		{
			name: "Get auto offset reset enviroment variable",
			args: args{
				name: "KAFKA_AUTO_OFFSET_RESET",
			},
			want: "101019",
		},
		{
			name: "Invalid - Get invalid enviroment variable",
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
	createMockEnvVars()
	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "one")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "two")
	defer deleteMockEnvVars()
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

			if got := RefreshEnviromentVariables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
