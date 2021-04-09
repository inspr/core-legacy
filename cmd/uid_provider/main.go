package main

import (
	"gitlab.inspr.dev/inspr/core/cmd/uid_provider/api"
	"gitlab.inspr.dev/inspr/core/cmd/uid_provider/client"
)

func main() {
	redisClient := client.NewRedisClient()
	api.Run(redisClient)
}
