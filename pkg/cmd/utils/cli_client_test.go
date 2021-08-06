package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/controller"
)

func TestGetCliClient(t *testing.T) {
	tests := []struct {
		name string
		want controller.Interface
	}{
		{
			name: "getCliClient working",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCliClient()
			if !reflect.DeepEqual(got, defaults.client) {
				t.Errorf(
					"GetCliClient() = %v, want %v",
					got,
					defaults.client,
				)
			}
		})
	}
}

func Test_setGlobalClient(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "setGlobalClient-working",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaults.client = nil
			setGlobalClient()
			if defaults.client == nil {
				t.Errorf(
					"GetCliClient() = %v, want %v",
					defaults.client,
					"non-nil",
				)
			}
		})
	}
}

func Test_setGlobalOutput(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "setGlobalOutput-working",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaults.out = nil
			setGlobalOutput()
			if defaults.out == nil {
				t.Errorf(
					"setGlobalOutput() = %v, want %v",
					defaults.out,
					"non-nil",
				)
			}
		})
	}
}

func TestGetCliOutput(t *testing.T) {
	tests := []struct {
		name string
		want io.Writer
	}{
		{
			name: "GetCliOutput-working",
			want: os.Stdout,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaults.out = nil
			got := GetCliOutput()
			if got != tt.want {
				t.Errorf(
					"GetCliOutput() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestSetOutput(t *testing.T) {
	tests := []struct {
		name    string
		wantOut string
	}{
		{
			name:    "testSetOutput-working",
			wantOut: "magic-test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			SetOutput(out)
			fmt.Fprintf(defaults.out, tt.wantOut)
			gotOut := out.String()

			if gotOut != tt.wantOut {
				t.Errorf(
					"SetOutput() = %v, want %v",
					gotOut,
					tt.wantOut,
				)
			}
		})
	}
}

func TestSetClient(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "setClient-working",
			args: args{
				url: "mock_url",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaults.client = nil
			SetClient(tt.args.url)
			if defaults.client == nil {
				t.Errorf("wanted non nil structure")
			}
		})
	}
}

func TestSetMockedClient(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "setMockedClient-working",
			args: args{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaults.client = nil
			SetMockedClient(tt.args.err)
			if defaults.client == nil {
				t.Errorf("wanted non nil structure")
			}
		})
	}
}
