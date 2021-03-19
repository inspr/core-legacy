package kafkasc

import (
	"os"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/environment"
)

func Test_fromTopicNonPRD(t *testing.T) {
	createMockEnv()
	defer deleteMockEnv()
	os.Setenv("INSPR_ENV", "test")
	os.Setenv("INSPR_APP_CTX", "random.app1")
	environment.RefreshEnviromentVariables()
	type args struct {
		topic string
	}
	tests := []struct {
		name string
		args args
		want messageChannel
	}{
		{
			name: "Non-PRD Environment topic",
			args: args{
				topic: "inspr-test-random.app1-nonPrdChan",
			},
			want: messageChannel{
				channel: "nonPrdChan",
				prefix:  "test",
				appCtx:  environment.GetInsprAppContext(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromTopic(tt.args.topic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromTopicPRD(t *testing.T) {
	createMockEnv()
	os.Setenv("INSPR_ENV", "")
	os.Setenv("INSPR_APP_CTX", "random.app1")
	defer deleteMockEnv()
	environment.RefreshEnviromentVariables()
	type args struct {
		topic string
	}
	tests := []struct {
		name string
		args args
		want messageChannel
	}{
		{
			name: "PRD Environment topic",
			args: args{
				topic: "inspr-random.app1-prdChan",
			},
			want: messageChannel{
				channel: "prdChan",
				prefix:  "",
				appCtx:  environment.GetInsprAppContext(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromTopic(tt.args.topic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toTopicNonPRD(t *testing.T) {
	createMockEnv()
	os.Setenv("INSPR_ENV", "test")
	os.Setenv("INSPR_APP_CTX", "random.app1")
	os.Setenv("nonPrdChan_RESOLVED", "random.app1.nonPrdChan")
	defer os.Unsetenv("INSPR_ENV")
	defer os.Unsetenv("INSPR_APP_CTX")
	defer os.Unsetenv("nonPrdChan_RESOLVED")

	defer deleteMockEnv()
	environment.RefreshEnviromentVariables()
	type args struct {
		channel messageChannel
		isPrd   bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "PRD Environment topic",
			args: args{
				channel: messageChannel{channel: "nonPrdChan", appCtx: "random.app1"},
				isPrd: func() bool {
					os.Unsetenv("INSPR_ENV")
					os.Setenv("INSPR_ENV", "test")
					return false
				}(),
			},
			want: "inspr-test-random.app1-nonPrdChan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.channel.toTopic(); got != tt.want {
				t.Errorf("toTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toTopicPRD(t *testing.T) {
	createMockEnv()
	os.Setenv("INSPR_ENV", "")
	os.Setenv("INSPR_APP_CTX", "random.app1")
	defer deleteMockEnv()
	environment.RefreshEnviromentVariables()
	type args struct {
		channel messageChannel
		isPrd   bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "PRD Environment topic",
			args: args{
				channel: messageChannel{channel: "prdChan", appCtx: "random.app1"},
				isPrd: func() bool {
					os.Unsetenv("INSPR_ENV")
					os.Setenv("INSPR_ENV", "")
					return true
				}(),
			},
			want: "inspr-random.app1-prdChan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.channel.toTopic(); got != tt.want {
				t.Errorf("toTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}
