package transports

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNewUnixSocketClient(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Valid unix socket client",
			args: args{
				addr: "/inspr/unixsockettest.sock",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUnixSocketClient(tt.args.addr)
			gotDial, err := NewUnixSocketClient(tt.args.addr).Transport.(*http.Transport).Dial("unix", tt.args.addr)

			if err != nil {
				fmt.Println(err)
				return
			}

			if got.Transport == nil && gotDial == nil {
				t.Errorf("Transport %v, Dial %v", got.Transport, gotDial)
			}
		})
	}
}
