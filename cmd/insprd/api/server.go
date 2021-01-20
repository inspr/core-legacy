package api

import ctrl "gitlab.inspr.dev/inspr/core/cmd/insprd/api/controllers"

var server ctrl.Server

// Run is the server start up function
func Run() {
	server.Init()
	server.Run(":8080")
}
