package main

import (
	"context"

	"inspr.dev/inspr/cmd/uidp/api"
	"inspr.dev/inspr/cmd/uidp/client"
)

func main() {
	ctx := context.Background()
	redisClient := client.NewRedisClient()
	api.Run(ctx, redisClient)
}
