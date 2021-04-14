package main

import (
	"context"

	"gitlab.inspr.dev/inspr/core/cmd/uid_provider/api"
	"gitlab.inspr.dev/inspr/core/cmd/uid_provider/client"
)

func main() {
	ctx := context.Background()
	redisClient := client.NewRedisClient()
	api.Run(redisClient, ctx)
}
