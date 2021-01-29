package api

import (
	ctrl "gitlab.inspr.dev/inspr/core/cmd/insprd/api/controllers"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
)

var server ctrl.Server

// Run is the server start up function
func Run() {
	// server.Init(mocks.MockMemoryManager(nil))
	server.Init(tree.GetTreeMemory())
	server.Run(":8080")
}
