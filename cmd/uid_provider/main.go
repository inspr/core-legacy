package main

import (
	"context"

	"github.com/inspr/inspr/cmd/uid_provider/api"
	"github.com/inspr/inspr/cmd/uid_provider/client"
)

func main() {
	ctx := context.Background()
	redisClient := client.NewRedisClient()
	api.Run(ctx, redisClient)
}
