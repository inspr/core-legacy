package main

import (
	"context"

	"inspr.dev/inspr/cmd/uid_provider/api"
	"inspr.dev/inspr/cmd/uid_provider/client"
)

func main() {
	ctx := context.Background()
	redisClient := client.NewRedisClient()
	api.Run(ctx, redisClient)
}
