package sidecarserv

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/sidecar_old/models"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Init(t *testing.T) {
	type args struct {
		r models.Reader
		w models.Writer
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Init(tt.args.r, tt.args.w)
		})
	}
}

func TestServer_Run(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Run(tt.args.ctx)
		})
	}
}

func Test_gracefulShutdown(t *testing.T) {
	type args struct {
		w   *http.Server
		r   *http.Server
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gracefulShutdown(tt.args.w, tt.args.r, tt.args.err)
		})
	}
}
