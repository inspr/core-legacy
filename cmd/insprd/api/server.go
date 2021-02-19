package api

import (
	ctrl "gitlab.inspr.dev/inspr/core/cmd/insprd/api/controllers"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
)

var server ctrl.Server

// Run is the server start up function
func Run(mm memory.Manager, nOp operators.NodeOperatorInterface, cOp operators.ChannelOperatorInterface) {
	// server.Init(mocks.MockMemoryManager(nil))
	server.Init(mm, nOp, cOp)
	server.Run(":8080")
}
