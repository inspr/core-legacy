package sidecarserv

import (
	"context"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

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
		{
			name: "sucessful_init",
			s:    &Server{},
			args: args{
				r: mockServer(nil).Reader,
				w: mockServer(nil).Writer,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Init(tt.args.r, tt.args.w)
			if _, err := tt.s.Reader.ReadMessage("channel"); err != nil {
				t.Errorf("function ReadMessage error. wanted %v got %v", nil, err)
			}
			if err := tt.s.Reader.CommitMessage("channel"); err != nil {
				t.Errorf("function CommitMessage error. wanted %v got %v", nil, err)
			}
			if err := tt.s.Writer.WriteMessage("channel", models.Message{}); err != nil {
				t.Errorf("function WriteMessage error. wanted %v got %v", nil, err)
			}
		})
	}
}

func TestServer_InitRoutes(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.InitRoutes()
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
