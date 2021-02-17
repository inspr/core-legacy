package transports

import (
	"net"
	"net/http"
)

// NewUnixSocketClient returns a client that conects via unix socket
func NewUnixSocketClient(addr string) http.Client {
	return http.Client{
		Transport: &http.Transport{
			Dial: func(string, string) (net.Conn, error) {
				return net.Dial("unix", addr)
			},
		},
	}
}
