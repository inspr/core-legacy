package transports

import (
	"net"
	"net/http"
)

// NewUnixSocketClient returna um client que conecta via unix socket
func NewUnixSocketClient(addr string) http.Client {
	sockAddr := "/inspr/" + addr + ".sock"
	return http.Client{
		Transport: &http.Transport{
			Dial: func(string, string) (net.Conn, error) {
				return net.Dial("unix", sockAddr)
			},
		},
	}
}
