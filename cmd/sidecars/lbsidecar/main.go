package main

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	kafkasc "inspr.dev/inspr/cmd/sidecars/kafka/client"
	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/sidecars/lbsidecar"
	"inspr.dev/inspr/pkg/sidecars/models"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = logs.Logger(zap.Fields(zap.String("section", "lb-sidecar-server")))
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

	kafkaHandler := models.NewBrokerHandler("kafka", reader, writer)

	logger.Info("initializing LB Sidecar server")
	lbServer := lbsidecar.Init(kafkaHandler)

	logger.Info("running LB Sidecar server")
	if err := lbServer.Run(ctx); err != nil {
		panic(err.Error())
	}
}
