package main

import (
	"context"

	"github.com/inspr/inspr/pkg/sidecars/lbsidecar"
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

	logger.Info("initializing LB Sidecar server")
	lbServer := lbsidecar.Init()

	logger.Info("running LB Sidecar server")
	if err := lbServer.Run(ctx); err != nil {
		panic(err.Error())
	}
}
