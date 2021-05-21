// package main

// import (
// 	"context"
// 	"fmt"

// 	kafkasc "github.com/inspr/inspr/cmd/sidecars/kafka/client"
// 	"github.com/inspr/inspr/pkg/environment"
// 	"github.com/inspr/inspr/pkg/sidecar_old/models"
// 	sidecarserv "github.com/inspr/inspr/pkg/sidecar_old/server"
// 	"go.uber.org/zap"
// )

// var logger *zap.Logger

// // init is called after all the variable declarations in the package have evaluated
// // their initializers, and those are evaluated only after all the imported packages
// // have been initialized
// func init() {
// 	logger, _ = zap.NewDevelopment(zap.Fields(zap.String("section", "kafka-sidecar-server")))
// }

// func main() {
// 	ctx := context.Background()
// 	var reader models.Reader
// 	var writer models.Writer
// 	var err error

// 	logger.Info("instantiating Kafka Sidecar reader")
// 	if len(environment.GetInputChannels()) != 0 {
// 		reader, err = kafkasc.NewReader()
// 		if err != nil {
// 			logger.Error("unable to instantiate Kafka Sidecar reader")

// 			fmt.Println(err)
// 			return
// 		}
// 	}

// 	logger.Info("instantiating Kafka Sidecar writer")
// 	if len(environment.GetOutputChannels()) != 0 {
// 		writer, err = kafkasc.NewWriter(false)
// 		if err != nil {
// 			logger.Error("unable to instantiate Kafka Sidecar writer")

// 			fmt.Println(err)
// 			return
// 		}
// 	}
// 	s := sidecarserv.NewServer()

// 	logger.Info("initializing Kafka Sidecar server")
// 	s.Init(reader, writer)

// 	logger.Info("running Kafka Sidecar server")
// 	s.Run(ctx)
// }
