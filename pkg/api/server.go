package api

import (
	"github.com/inspr/inspr/cmd/insprd/memory"
	ctrl "github.com/inspr/inspr/pkg/api/controllers"
	"github.com/inspr/inspr/pkg/auth"
)

var server ctrl.Server

// Run is the server start up function
func Run(mm memory.Manager, auth auth.Auth) {
	// server.Init(mocks.MockMemoryManager(nil))
	server.Init(mm, auth)
	server.Run(":8080")
}
