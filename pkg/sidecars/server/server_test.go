package sidecarserv

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/sidecars/models"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
		want *Server
	}{
		{
			name: "test_basic_server_creation",
			want: &Server{},
		},
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
	test := struct {
		name    string
		addr    string
		channel string
	}{
		name:    "basic init test",
		addr:    "localhost:8080",
		channel: "testing",
	}

	createMockEnvVars() // creates mock values for test
	defer deleteMockEnvVars()

	t.Run("basic init test", func(t *testing.T) {
		s := &Server{
			writeAddr: test.addr,
		}
		r := mockReader{
			readMessage: func(ctx context.Context, channel string) (models.BrokerMessage, error) {
				return models.BrokerMessage{}, nil
			},
			commit: func(ctx context.Context, channel string) error {
				return nil
			},
		}
		w := mockWriter{
			writeMessage: func(channel string, message interface{}) error {
				return nil
			},
		}
		s.Init(r, w, models.ConnectionVariables{
			ReadEnvVar:  "INSPR_SIDECAR_READ_PORT",
			WriteEnvVar: "INSPR_SIDECAR_WRITE_PORT",
		})

		// checking reader methods
		if got := s.Reader.Commit(context.Background(), test.channel); got != nil {
			t.Errorf("expected CommitMessage() == nil, received %v", got)
		}
		if _, got := s.Reader.ReadMessage(context.Background(), test.channel); got != nil {
			t.Errorf("expected CommitMessage() == nil, received %v", got)
		}
		if got := s.Writer.WriteMessage("channel", "msg"); got != nil {
			t.Errorf("expected CommitMessage() == nil, received %v", got)
		}
	})
}

func TestServer_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	environment.SetMockEnv()
	server := &Server{
		writeAddr: ":3001",
		Reader: mockReader{
			readMessage: func(ctx context.Context, channel string) (models.BrokerMessage, error) {
				<-ctx.Done()
				return models.BrokerMessage{}, ctx.Err()
			},
			commit: func(ctx context.Context, channel string) error {
				return nil
			},
		},
	}
	done := make(chan struct{})
	go func() { server.Run(ctx); done <- struct{}{} }()

	time.Sleep(time.Second)
	if !server.runningRead {
		t.Error("Server_Run read message not initialized")
	}
	if !server.runningWrite {
		t.Error("Server_Run write message not initialized")
	}

	deadContext, cancelDead := context.WithTimeout(context.Background(), time.Second)
	defer cancelDead()
	cancel()
	select {
	case <-done:

		if server.runningRead {
			t.Error("Server_Run read message not finalized")

		}
		if server.runningWrite {
			t.Error("Server_Run write message not finalized")
		}
	case <-deadContext.Done():
		t.Errorf("Server_Run timeout")
	}

}
