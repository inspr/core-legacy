package api

import (
	"context"

	"github.com/inspr/inspr/cmd/uid_provider/api/controller"
	"github.com/inspr/inspr/cmd/uid_provider/client"
)

var server controller.Server

// Run runs the UID Provider API server
func Run(rdb client.RedisManager, ctx context.Context) {
	server.Init(rdb, ctx)
	server.Run(":9001")
}
