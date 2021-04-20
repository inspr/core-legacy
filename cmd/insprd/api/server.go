package api

import (
	ctrl "github.com/inspr/inspr/cmd/insprd/api/controllers"
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/operators"
	"github.com/inspr/inspr/pkg/auth"
)

var server ctrl.Server

// Run is the server start up function
func Run(mm memory.Manager, op operators.OperatorInterface, auth auth.Auth) {
	// server.Init(mocks.MockMemoryManager(nil))
	server.Init(mm, op, auth)
	server.Run(":8080")
}
