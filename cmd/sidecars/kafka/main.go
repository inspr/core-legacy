package main

import (
	"context"
	"fmt"

	kafkasc "github.com/inspr/inspr/cmd/sidecars/kafka/client"
	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/meta/brokers"

	"github.com/inspr/inspr/pkg/sidecars/models"
	sidecarserv "github.com/inspr/inspr/pkg/sidecars/server"
	"go.uber.org/zap"
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
		reader, err = kafkasc.NewReader() // alterar metodo para comply a nova interface
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
