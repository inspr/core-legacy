package main

import (
	"context"
	"fmt"

	kafkasc "gitlab.inspr.dev/inspr/core/cmd/sidecars/kafka/client"
	sidecarserv "gitlab.inspr.dev/inspr/core/pkg/sidecar/server"
)

func main() {
	ctx := context.Background()
	reader, err := kafkasc.NewReader()
	if err != nil {
		fmt.Println(err)
		return
	}

	writer, err := kafkasc.NewWriter(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	s := sidecarserv.NewServer()
	s.Init(reader, writer)

	s.Run(ctx)
}
