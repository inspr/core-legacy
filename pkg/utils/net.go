package utils

import "net"

// GetFreePorts asks the kernel for free open ports that are ready to use
func GetFreePorts(count int) ([]int, error) {
	var ports []int
	for i := 0; i < count; i++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return nil, err
		}

		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, err
		}
		defer listener.Close()
		ports = append(ports, listener.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}
