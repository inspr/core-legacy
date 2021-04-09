package api

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	ctrl "gitlab.inspr.dev/inspr/core/pkg/api/controllers"
)

var server ctrl.Server

// Run is the server start up function
func Run(mm memory.Manager, op operators.OperatorInterface) {
	// server.Init(mocks.MockMemoryManager(nil))
	server.Init(mm, op)
	server.Run(":8080")
}
