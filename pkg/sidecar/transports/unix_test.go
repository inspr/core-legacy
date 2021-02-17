package transports

import (
	"net"
	"net/http"
	"testing"
)

func mockHTTPClient(addr string) http.Client {
	return http.Client{
		Transport: &http.Transport{
			Dial: func(string, string) (net.Conn, error) {
				return net.Dial("unix", addr)
			},
		},
	}
}

func TestNewUnixSocketClient(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want http.Client
	}{
		{
			name: "Valid unix socket client",
			args: args{
				addr: "/inspr/unixsockettest.sock",
			},
			want: mockHTTPClient("/inspr/unixsockettest.sock"),
		},
		{
			name: "Invalid unix socket client",
			args: args{
				addr: "/inspr/unixsockettest.sock",
			},
			want: mockHTTPClient("/inspr/unixsockettest.sock"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUnixSocketClient(tt.args.addr)
			if got.Transport == tt.want.Transport {
				t.Errorf("NewUnixSocketClient().Transport = %v, want %v", got.Transport, tt.want.Transport)
			}
		})
	}
}
