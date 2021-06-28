package main

import (
	"context"
	"fmt"

	kafkasc "inspr.dev/inspr/cmd/sidecars/kafka/client"
	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/meta/brokers"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/sidecars/models"
	sidecarserv "inspr.dev/inspr/pkg/sidecars/server"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "kafka-sidecar-server")))
}

func main() {
	ctx := context.Background()
	var reader models.Reader
	var writer models.Writer
	var err error

	logger.Info("instantiating Kafka Sidecar reader")
	if len(environment.GetInputChannelsData()) != 0 {
		reader, err = kafkasc.NewReader()
		if err != nil {
			logger.Error("unable to instantiate Kafka Sidecar reader")

			fmt.Println(err)
			return
		}
	}

	logger.Info("instantiating Kafka Sidecar writer")
	if len(environment.GetOutputChannelsData()) != 0 {
		writer, err = kafkasc.NewWriter()
		if err != nil {
			logger.Error("unable to instantiate Kafka Sidecar writer")

			fmt.Println(err)
			return
		}
	}

	logger.Info("initializing Kafka Sidecar server")
	s := sidecarserv.Init(reader, writer, brokers.Kafka)

	logger.Info("running Kafka Sidecar server")
	err = s.Run(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
