package api

import ctrl "github.com/inspr/inspr/cmd/authsvc/api/controllers"

var server ctrl.Server

// Run is the server start up function
func Run() {
	server.Init()
	server.Run(":8081")
}
