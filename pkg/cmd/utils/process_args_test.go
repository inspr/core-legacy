package utils

import (
	"testing"

	"github.com/inspr/inspr/pkg/ierrors"
)

func TestCheckEmptyArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]string
		wantErr     bool
		wantWrapper bool
	}{
		{
			name: "no_errors",
			args: map[string]string{
				"test": "value",
			},
			wantErr: false,
		},
		{
			name: "one_error",
			args: map[string]string{
				"test": "",
			},
			wantErr: true,
		},
		{
			name: "multiple_errors",
			args: map[string]string{
				"A": "",
				"B": "",
				"C": "",
			},
			wantErr:     true,
			wantWrapper: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckEmptyArgs(tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckEmptyArgs() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantWrapper && ierrors.Unwrap(err) == nil {
				t.Errorf("CheckEmptyArgs() multiple errors, wanted wrapper")
			}
		})
	}
}

func TestProcessArg(t *testing.T) {
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
			got, got1, err := ProcessArg(tt.args.arg, tt.args.scope)
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

func TestProcessAliasArg(t *testing.T) {
	type args struct {
		arg   string
		scope string
	}
	tests := []struct {
		name    string
		args    args
		path    string
		alias   string
		wantErr bool
	}{
		{
			name: "Arg is a invalid alias structure name - it should return a error",
			args: args{
				arg:   "invalid!name",
				scope: "",
			},
			wantErr: true,
		},
		{
			name: "Arg is a valid alias structure name",
			args: args{
				arg:   "helloWorld.alias",
				scope: "app1",
			},
			path:    "app1",
			alias:   "helloWorld.alias",
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
				arg:   "hello.World.alias",
				scope: "app1",
			},
			path:    "app1.hello",
			alias:   "World.alias",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ProcessAliasArg(tt.args.arg, tt.args.scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessAliasArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.path {
				t.Errorf("ProcessAliasArg() got = %v, want %v", got, tt.path)
			}
			if got1 != tt.alias {
				t.Errorf("ProcessAliasArg() got1 = %v, want %v", got1, tt.alias)
			}
		})
	}
}
