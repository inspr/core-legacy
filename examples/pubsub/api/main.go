package main

import (
	controller "gitlab.inspr.dev/inspr/core/examples/pubsub/api/controller"
)

var server controller.Server

// main is the server start up function
func main() {
	server.Init()
	server.Run(":8080")
}
