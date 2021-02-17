package api

import (
	ctrl "gitlab.inspr.dev/inspr/core/cmd/insprd/api/controllers"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
)

var server ctrl.Server

// Run is the server start up function
func Run(mm memory.Manager) {
	// server.Init(mocks.MockMemoryManager(nil))
	server.Init(mm)
	server.Run(":8080")
}
