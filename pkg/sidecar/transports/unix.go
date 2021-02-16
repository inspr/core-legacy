package transports

import (
	"net"
	"net/http"
)

// NewUnixSocketClient returna um client que conecta via unix socket
func NewUnixSocketClient(addr string) http.Client {
	return http.Client{
		Transport: &http.Transport{
			Dial: func(string, string) (net.Conn, error) {
				return net.Dial("unix", addr)
			},
		},
	}
}
