package main

import (
	"github.com/inspr/inspr/cmd/uid_provider/api"
	"github.com/inspr/inspr/cmd/uid_provider/client"
)

func main() {
	redisClient := client.NewRedisClient()
	api.Run(redisClient)
}
