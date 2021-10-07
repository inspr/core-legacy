package lbsidecar

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// The environment variables used here are declared in handlers_test.go

func TestInit(t *testing.T) {
	createMockEnvVars()
	defer deleteMockEnvVars()

	tests := []struct {
		name string
		want *Server
	}{
		{
			name: "Initializes a new server",
			want: &Server{
				writeAddr:     fmt.Sprintf(":%s", os.Getenv("INSPR_LBSIDECAR_WRITE_PORT")),
				readAddr:      fmt.Sprintf(":%s", os.Getenv("INSPR_LBSIDECAR_READ_PORT")),
				channelMetric: make(map[string]channelMetric),
				routeMetric:   make(map[string]routeMetric),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Run(t *testing.T) {
	createMockEnvVars()
	defer deleteMockEnvVars()

	type fields struct {
		writeAddr string
		readAddr  string
	}
	tests := []struct {
		name         string
		fields       fields
		requestWrite bool
		requestRead  bool
	}{
		{
			name: "Runs server and cancel its context afterwards",
			fields: fields{
				writeAddr: fmt.Sprintf(":%s", os.Getenv("INSPR_LBSIDECAR_WRITE_PORT")),
				readAddr:  fmt.Sprintf(":%s", os.Getenv("INSPR_LBSIDECAR_READ_PORT")),
			},
		},
		{
			name: "Tries to create read server on already-used port",
			fields: fields{
				writeAddr: fmt.Sprintf(":%s", os.Getenv("INSPR_LBSIDECAR_WRITE_PORT")),
				readAddr:  fmt.Sprintf(":%s", os.Getenv("INSPR_LBSIDECAR_READ_PORT")),
			},
			requestRead: true,
		},
		{
			name: "Tries to create read server on already-used port",
			fields: fields{
				writeAddr: fmt.Sprintf(":%s", os.Getenv("INSPR_LBSIDECAR_WRITE_PORT")),
				readAddr:  fmt.Sprintf(":%s", os.Getenv("INSPR_LBSIDECAR_READ_PORT")),
			},
			requestWrite: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				writeAddr: tt.fields.writeAddr,
				readAddr:  tt.fields.readAddr,
			}

			if tt.requestRead && tt.requestWrite {
				t.Error("for testing purposes, choose only one of 'requestRead' and 'requestWrite'")
				return
			} else if !tt.requestRead && !tt.requestWrite {
				errChan := make(chan error)
				ctx, cancel := context.WithCancel(context.Background())

				go func() { errChan <- s.Run(ctx) }()
				cancel()

				if err := <-errChan; err.Error() != "context canceled" {
					t.Errorf("expected 'context canceled', got '%v'", err)
				}
			} else if tt.requestRead {
				port := strings.Split(tt.fields.readAddr, ":")[1]
				auxServer := createMockedServer(port, "randCh", "randMsg")
				auxServer.Start()
				defer auxServer.Close()

				if err := s.Run(context.Background()); err.Error() != "listen tcp :1137: bind: address already in use" {
					t.Errorf("expected bind address error, got '%v'", err.Error())
				}
			} else {
				port := strings.Split(tt.fields.writeAddr, ":")[1]
				auxServer := createMockedServer(port, "randCh", "randMsg")
				auxServer.Start()
				defer auxServer.Close()

				if err := s.Run(context.Background()); err.Error() != "listen tcp :1127: bind: address already in use" {
					t.Errorf("expected bind address error, got '%v'", err.Error())
				}
			}
		})
	}
}
