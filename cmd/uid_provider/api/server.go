package api

import (
	"context"

	"github.com/inspr/inspr/cmd/uid_provider/api/controller"
	"github.com/inspr/inspr/cmd/uid_provider/client"
)

var server controller.Server

// Run runs the UID Provider API server
func Run(ctx context.Context, rdb client.RedisManager) {
	server.Init(ctx, rdb)
	server.Run(":9001")
}
