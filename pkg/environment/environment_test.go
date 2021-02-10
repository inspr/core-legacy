package environment

import (
	"os"
	"reflect"
	"testing"
)

// createMockEnvVars - sets up the env values to be used in the tests functions
// createMockEnvVars - sets up the env values to be used in the tests functions
func createMockEnvVars() {
	os.Setenv("INSPR_INPUT_CHANNELS", "inp1;inp2;inp3")
	os.Setenv("INSPR_OUTPUT_CHANNELS", "out1;out2;out3")
	os.Setenv("INSPR_UNIX_SOCKET", "/addr/to/socket")
	os.Setenv("INSPR_APP_CTX", "random.ctx")
	os.Setenv("INSPR_ENV", "test")
}

// deleteMockEnvVars - deletes the env values used in the tests functions
func deleteMockEnvVars() {
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("INSPR_UNIX_SOCKET")
	os.Unsetenv("INSPR_APP_CTX")
	os.Unsetenv("INSPR_ENV")
}

func mockInsprEnvironment() *InsprEnvironmentVariables {
	return &InsprEnvironmentVariables{
		InputChannels:    "inp1;inp2;inp3",
		OutputChannels:   "out1;out2;out3",
		UnixSocketAddr:   "/addr/to/socket",
		InsprAppContext:  "random.ctx",
		InsprEnvironment: "test",
	}
}

func TestGetEnvironment(t *testing.T) {
	createMockEnvVars()
	defer deleteMockEnvVars()
	tests := []struct {
		name string
		want *InsprEnvironmentVariables
	}{
		{
			name: "Get all environment variables",
			want: mockInsprEnvironment(),
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
	createMockEnvVars()
	defer deleteMockEnvVars()
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get input channel enviroment variable",
			args: args{
				name: "INSPR_INPUT_CHANNELS",
			},
			want: "inp1;inp2;inp3",
		},
		{
			name: "Get output channel enviroment variable",
			args: args{
				name: "INSPR_OUTPUT_CHANNELS",
			},
			want: "out1;out2;out3",
		},
		{
			name: "Get unix socket enviroment variable",
			args: args{
				name: "INSPR_UNIX_SOCKET",
			},
			want: "/addr/to/socket",
		},
		{
			name: "Invalid - Get invalid enviroment variable",
			args: args{
				name: "INSPR_INVALID_ENV_VAR",
			},
			want: "[ENV VAR] INSPR_INVALID_ENV_VAR not found",
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
	os.Setenv("INSPR_INPUT_CHANNELS", "one")
	os.Setenv("INSPR_OUTPUT_CHANNELS", "two")
	os.Setenv("INSPR_UNIX_SOCKET", "three")
	os.Setenv("INSPR_APP_CTX", "four")
	os.Setenv("INSPR_ENV", "five")
	defer deleteMockEnvVars()
	tests := []struct {
		name    string
		refresh bool
		want    *InsprEnvironmentVariables
	}{
		{
			name:    "Changed and refreshed environment variables",
			refresh: true,
			want: &InsprEnvironmentVariables{
				InputChannels:    "one",
				OutputChannels:   "two",
				UnixSocketAddr:   "three",
				InsprAppContext:  "four",
				InsprEnvironment: "five",
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
