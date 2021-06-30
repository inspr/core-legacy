package api

import (
	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/operators"
	ctrl "inspr.dev/inspr/pkg/api/controllers"
	"inspr.dev/inspr/pkg/auth"
)

var server ctrl.Server

// Run is the server start up function
func Run(mem memory.Manager, op operators.OperatorInterface, auth auth.Auth) {
	// server.Init(mocks.MockMemoryManager(nil))
	server.Init(mem, op, auth)
	server.Run(":8080")
}
