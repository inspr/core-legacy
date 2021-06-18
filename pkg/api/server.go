package api

import (
	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/cmd/insprd/operators"
	ctrl "inspr.dev/inspr/pkg/api/controllers"
	"inspr.dev/inspr/pkg/auth"
)

var server ctrl.Server

// Run is the server start up function
func Run(mm memory.Manager, op operators.OperatorInterface, auth auth.Auth, bm brokers.Manager) {
	// server.Init(mocks.MockMemoryManager(nil))
	server.Init(mm, op, auth, bm)
	server.Run(":8080")
}
