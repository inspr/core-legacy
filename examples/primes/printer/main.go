package main

import (
	"context"
	"fmt"
	"log"
	"time"

	dappclient "gitlab.inspr.dev/inspr/core/pkg/client"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

func main() {
	// sets up client for sidecar
	c := dappclient.NewAppClient()

	// sets up ticker
	ticker := time.NewTicker(2 * time.Second)
	ctx := context.Background()
	chName := "primes_ch2"

	for {
		select {
		case <-ticker.C:
			var message models.Message
			err := c.ReadMessage(ctx, chName, &message)
			if err != nil {
				log.Println(err.Error())
			}
			fmt.Println("Message Content -> ", message.Data)

			err = c.CommitMessage(ctx, chName)
			if err != nil {
				log.Println(err.Error())
			}

		}
	}
}
