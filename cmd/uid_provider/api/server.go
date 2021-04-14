package api

import (
	ctrl "github.com/inspr/inspr/cmd/uid_provider/api/controllers"
	"github.com/inspr/inspr/cmd/uid_provider/client"
)

var server ctrl.Server

// Run runs the UID Provider API server
func Run(rdb client.RedisManager) {
	server.Init(rdb)
	server.Run(":9001")
}
