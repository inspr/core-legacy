package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/controller"
	"inspr.dev/inspr/pkg/ierrors"
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

func TestRequestErrorMessage(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ierror-unauthorized",
			args: args{
				ierrors.New("").Unauthorized(),
			},
			want: "failed to authenticate with the cluster. Is your token configured correctly?\n",
		},
		{
			name: "ierror-forbidden",
			args: args{
				ierrors.New("").Forbidden(),
			},
			want: "forbidden operation, please check for the scope.\n",
		},
		{
			name: "ierror-unknown",
			args: args{
				err: errors.New("mock-error"),
			},
			want: ierrors.Wrap(
				ierrors.New("mock-error"),
				"unknown inspr error",
			).Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			RequestErrorMessage(tt.args.err, w)
			if gotW := w.String(); gotW != tt.want {
				t.Errorf("RequestErrorMessage() = %v, want %v", gotW, tt.want)
			}
		})
	}
}
