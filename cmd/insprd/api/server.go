package api

import (
	ctrl "gitlab.inspr.dev/inspr/core/cmd/insprd/api/controllers"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/mocks"
)

var server ctrl.Server

// Run is the server start up function
func Run() {
	server.Init(mocks.MockMemoryManager(nil))
	server.Run(":8080")
}
