package api

import (
	"context"

	"inspr.dev/inspr/cmd/uidp/api/controller"
	"inspr.dev/inspr/cmd/uidp/client"
)

var server controller.Server

// Run runs the UID Provider API server
func Run(ctx context.Context, rdb client.RedisManager) {
	server.Init(ctx, rdb)
	server.Run(":9001")
}
