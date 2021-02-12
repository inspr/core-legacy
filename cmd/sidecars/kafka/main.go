package main

import (
	"context"

	kafkasc "gitlab.inspr.dev/inspr/core/cmd/sidecars/kafka/client"
	sidecarserv "gitlab.inspr.dev/inspr/core/pkg/sidecar/server"
)

func main() {
	ctx := context.Background()
	reader, err := kafkasc.NewReader()
	if err != nil {

	}

	writer, err := kafkasc.NewWriter(false)
	if err != nil {

	}

	s := sidecarserv.NewServer()
	s.Init(reader, writer)

	s.Run(ctx)
}
