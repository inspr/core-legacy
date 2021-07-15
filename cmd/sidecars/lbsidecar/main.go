package main

import (
	"context"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/sidecars/lbsidecar"
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

	logger.Info("initializing LB Sidecar server")
	lbServer := lbsidecar.Init()

	logger.Info("running LB Sidecar server")
	if err := lbServer.Run(ctx); err != nil {
		panic(err.Error())
	}
}
