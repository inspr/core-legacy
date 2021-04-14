package api

import (
	"context"

	api "gitlab.inspr.dev/inspr/core/cmd/uid_provider/api/controllers"
	"gitlab.inspr.dev/inspr/core/cmd/uid_provider/client"
)

var server api.Server

// Run runs the UID Provider API server
func Run(rdb client.RedisManager, ctx context.Context) {
	server.Init(rdb, ctx)
	server.Run(":9001")
}
