package sidecarserv

import (
	"context"
	"testing"
	"time"

	"github.com/inspr/inspr/cmd/insprd/memory/brokers"
	"github.com/inspr/inspr/pkg/environment"
)

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

		r := &mockReader{
			readMessage: func(ctx context.Context, channel string) ([]byte, error) {
				return nil, nil
			},
			commit: func(ctx context.Context, channel string) error {
				return nil
			},
		}
		w := &mockWriter{
			writeMessage: func(channel string, message []byte) error {
				return nil
			},
		}
		s := Init(r, w, brokers.Kafka)

		// checking reader methods
		if got := s.Reader.Commit(context.Background(), test.channel); got != nil {
			t.Errorf("expected CommitMessage() == nil, received %v", got)
		}
		if _, got := s.Reader.ReadMessage(context.Background(), test.channel); got != nil {
			t.Errorf("expected CommitMessage() == nil, received %v", got)
		}
		if got := s.Writer.WriteMessage("channel", []byte("msg*")); got != nil {
			t.Errorf("expected CommitMessage() == nil, received %v", got)
		}
	})
}

func TestServer_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	environment.SetMockEnv()
	server := &Server{
		inAddr: ":3001",
		Reader: &mockReader{
			readMessage: func(ctx context.Context, channel string) ([]byte, error) {
				<-ctx.Done()
				return nil, ctx.Err()
			},
			commit: func(ctx context.Context, channel string) error {
				return nil
			},
		},
		Writer: &mockWriter{},
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
