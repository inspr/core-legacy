package environment

import (
<<<<<<< HEAD
	"os"
=======
>>>>>>> story/core-42
	"reflect"
	"testing"
)

<<<<<<< HEAD
func mockInsprEnvironment() *InsprEnvironment {
	return &InsprEnvironment{
		InputChannels:  "inp1;inp2;inp3",
		OutputChannels: "out1;out2;out3",
		UnixSocketAddr: "/addr/to/socket",
	}
}

=======
>>>>>>> story/core-42
func TestGetEnvironment(t *testing.T) {
	tests := []struct {
		name string
		want *InsprEnvironment
	}{
<<<<<<< HEAD
		{
			name: "Get all environment variables",
			want: mockInsprEnvironment(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("INSPR_INPUT_CHANNELS", "inp1;inp2;inp3")
			os.Setenv("INSPR_OUTPUT_CHANNELS", "out1;out2;out3")
			os.Setenv("INSPR_UNIX_SOCKET", "/addr/to/socket")

			defer func() {
				recover()
			}()

=======
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
>>>>>>> story/core-42
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
<<<<<<< HEAD
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
			os.Setenv("INSPR_INPUT_CHANNELS", "inp1;inp2;inp3")
			os.Setenv("INSPR_OUTPUT_CHANNELS", "out1;out2;out3")
			os.Setenv("INSPR_UNIX_SOCKET", "/addr/to/socket")

			defer func() {
				recover()
			}()

=======
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
>>>>>>> story/core-42
			if got := getEnv(tt.args.name); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
