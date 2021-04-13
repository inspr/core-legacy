package main

import (
	controller "github.com/inspr/inspr/examples/pubsub/api/controller"
)

var server controller.Server

// main is the server start up function
func main() {
	server.Init()
	server.Run(":8080")
}
