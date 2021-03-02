package cli

import (
	"testing"

	"github.com/spf13/cobra"
	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"
	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
)

func TestNewDescribeCmd(t *testing.T) {
	tests := []struct {
		name          string
		checkFunctiom func(t *testing.T, got *cobra.Command)
	}{
		{
			name: "It should create a new describe command",
			checkFunctiom: func(t *testing.T, got *cobra.Command) {
				if got == nil {
					t.Errorf("NewDescribeCmd() not created successfully")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDescribeCmd()
			if tt.checkFunctiom != nil {
				tt.checkFunctiom(t, got)
			}
		})
	}
}

func Test_getClient(t *testing.T) {
	tests := []struct {
		name          string
		checkFunction func(t *testing.T, got *client.Client)
	}{
		{
			name: "It should return a controller client",
			checkFunction: func(t *testing.T, got *client.Client) {
				if got == nil {
					t.Errorf("getClient() = nil")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cliutils.GetClient()
			if tt.checkFunction != nil {
				tt.checkFunction(t, got)
			}
		})
	}
}

func Test_getScope(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "It should return the scope",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cliutils.GetScope()
			if (err != nil) != tt.wantErr {
				t.Errorf("getScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getScope() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processArg(t *testing.T) {
	type args struct {
		arg   string
		scope string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "Arg is a invalid structure name - it should return a error",
			args: args{
				arg:   "invalid!name",
				scope: "",
			},
			wantErr: true,
		},
		{
			name: "Arg is a valid structure name",
			args: args{
				arg:   "helloWorld",
				scope: "app1",
			},
			want:    "app1",
			want1:   "helloWorld",
			wantErr: false,
		},
		{
			name: "Arg is a invalid scope - it should return a error",
			args: args{
				arg:   "hello..World",
				scope: "app1",
			},
			wantErr: true,
		},
		{
			name: "Arg is a valid scope",
			args: args{
				arg:   "hello.World",
				scope: "app1",
			},
			want:    "app1.hello",
			want1:   "World",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := cliutils.ProcessArg(tt.args.arg, tt.args.scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("processArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("processArg() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("processArg() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
