package sidecarserv

import (
	"context"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"gitlab.inspr.dev/inspr/core/pkg/sidecar/transports"
)

func TestServer_Run(t *testing.T) {
	routes := []string{"commit", "writeMessage", "readMessage"}
	for _, r := range routes {

		t.Run("run_test/"+r, func(t *testing.T) {
			// ENV variables
			os.Setenv("INSPR_INPUT_CHANNELS", "")
			os.Setenv("INSPR_OUTPUT_CHANNELS", "")
			os.Setenv("UNIX_SOCKET_ADDRESS", unixSocketAddr)

			// SERVER
			var wg sync.WaitGroup
			wg.Add(1)
			defer wg.Wait()

			s := MockServer(nil)
			s.Init(s.Reader, s.Writer)
			go func() {
				s.Run(context.Background())
			}()

			go func() {
				defer wg.Done()
				time.Sleep(500 * time.Microsecond)

				c := transports.NewUnixSocketClient(unixSocketAddr)

				resp, err := c.Post("http://unix/"+r, "", nil)
				if err != nil {
					t.Errorf("Failed to make post to route '/commit'")
				}
				if resp.StatusCode != http.StatusBadRequest {
					t.Errorf("route '/commit' = %v, want %v", resp.StatusCode, http.StatusBadRequest)
				}
			}()
		})
	}
}

func TestServer_Cancel(t *testing.T) {
	t.Run("run_test/timeout", func(t *testing.T) {
		// ENV variables
		os.Setenv("INSPR_INPUT_CHANNELS", "")
		os.Setenv("INSPR_OUTPUT_CHANNELS", "")
		os.Setenv("UNIX_SOCKET_ADDRESS", unixSocketAddr)

		// SERVER
		var wg sync.WaitGroup
		wg.Add(1)
		defer wg.Wait()

		ctx, cancel := context.WithCancel(context.Background())

		s := MockServer(nil)
		s.Init(s.Reader, s.Writer)
		go func() {
			s.Run(ctx)
		}()

		go func() {
			defer wg.Done()
			time.Sleep(500 * time.Microsecond)
			cancel()
			time.Sleep(500 * time.Microsecond)

			c := transports.NewUnixSocketClient(unixSocketAddr)

			_, err := c.Post("http://unix/commit", "", nil)
			if err == nil {
				t.Errorf("Server should be down")
			}
		}()
	})
}
