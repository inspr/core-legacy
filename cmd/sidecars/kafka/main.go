package main

import (
	"context"
	"fmt"

	kafkasc "gitlab.inspr.dev/inspr/core/cmd/sidecars/kafka/client"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
	sidecarserv "gitlab.inspr.dev/inspr/core/pkg/sidecar/server"
)

func main() {
	ctx := context.Background()
	var reader models.Reader
	var writer models.Writer
	var err error
	if len(environment.GetInputChannels()) != 0 {
		reader, err = kafkasc.NewReader()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if len(environment.GetOutputChannels()) != 0 {
		writer, err = kafkasc.NewWriter(false)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	s := sidecarserv.NewServer()
	s.Init(reader, writer)

	s.Run(ctx)
}
