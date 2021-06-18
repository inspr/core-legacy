package api

import ctrl "inspr.dev/inspr/cmd/authsvc/api/controllers"

var server ctrl.Server

// Run is the server start up function
func Run() {
	server.Init()
	server.Run(":8081")
}
