package sidecarserv

import (
	"context"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	env "inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/sidecar/models"
	"inspr.dev/inspr/pkg/sidecar/transports"
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
	type args struct {
		r models.Reader
		w models.Writer
	}
	test := struct {
		name    string
		addr    string
		channel string
		args    args
	}{
		name: "basic init test",
		addr: "localhost:8080",
		args: args{
			r: MockServer(nil).Reader,
			w: &mockWriter{},
		},
		channel: "testing",
	}

	createMockEnvVars() // creates mock values for test
	defer deleteMockEnvVars()

	t.Run("basic init test", func(t *testing.T) {
		s := &Server{
			Mux:  MockServer(nil).Mux,
			addr: test.addr,
		}
		s.Init(test.args.r, test.args.w)

		// checking reader methods
		if got := s.Reader.Commit(test.channel); got != nil {
			t.Errorf("expected CommitMessage() == nil, received %v", got)
		}
		if _, got := s.Reader.ReadMessage(test.channel); got != nil {
			t.Errorf("expected CommitMessage() == nil, received %v", got)
		}
		if got := s.Writer.WriteMessage("channel", "msg"); got != nil {
			t.Errorf("expected CommitMessage() == nil, received %v", got)
		}
	})
}

func TestServer_Run(t *testing.T) {
	routes := []string{"commit", "writeMessage", "readMessage"}
	env.SetMockEnv()

	for _, r := range routes {
		t.Run("run_test/"+r, func(t *testing.T) {

			// SERVER
			var wg sync.WaitGroup
			wg.Add(1)
			defer wg.Wait()

			s := MockServer(nil)
			s.Init(s.Reader, s.Writer)
			s.addr = "./test.sock"

			go func() {
				s.Run(context.Background())
			}()

			go func() {
				defer wg.Done()
				time.Sleep(500 * time.Microsecond)

				// env mock socket addr
				c := transports.NewUnixSocketClient("./test.sock")

				resp, err := c.Post("http://unix/"+r, "", nil)
				if err != nil {
					t.Errorf("Failed to make post to route '/commit'")
					return
				}
				if resp.StatusCode != http.StatusBadRequest {
					t.Errorf("route '/commit' = %v, want %v", resp.StatusCode, http.StatusBadRequest)
				}
			}()
		})
	}
	env.UnsetMockEnv()
}

func TestServer_Cancel(t *testing.T) {
	env.SetMockEnv()

	t.Run("run_test/timeout", func(t *testing.T) {
		// SERVER
		var wg sync.WaitGroup
		wg.Add(1)
		defer wg.Wait()

		ctx, cancel := context.WithCancel(context.Background())

		s := MockServer(nil)
		s.Init(s.Reader, s.Writer)
		s.addr = "./test.sock"
		go func() {
			s.Run(ctx)
		}()

		go func() {
			defer wg.Done()
			time.Sleep(500 * time.Microsecond)
			cancel()
			time.Sleep(500 * time.Microsecond)

			c := transports.NewUnixSocketClient(env.GetUnixSocketAddress())

			_, err := c.Post("http://unix/commit", "", nil)
			if err == nil {
				t.Errorf("Server should be down")
			}
		}()
	})
	env.UnsetMockEnv()
}
