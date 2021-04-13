package api

import (
	ctrl "github.com/inspr/inspr/cmd/insprd/api/controllers"
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/operators"
)

var server ctrl.Server

// Run is the server start up function
func Run(mm memory.Manager, op operators.OperatorInterface) {
	// server.Init(mocks.MockMemoryManager(nil))
	server.Init(mm, op)
	server.Run(":8080")
}
