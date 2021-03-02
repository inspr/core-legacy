package sidecarserv

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	env "gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/transports"
)

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
