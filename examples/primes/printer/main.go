package main

import (
	"context"
	"fmt"
	"log"
	"time"

	dappclient "inspr.dev/inspr/pkg/client"
)

func main() {
	// sets up client for sidecar
	c := dappclient.NewAppClient()

	// sets up ticker
	ticker := time.NewTicker(2 * time.Second)
	ctx := context.Background()
	chName := "input"
	fmt.Println("starting...")
	for range ticker.C {
		var message struct {
			Message struct {
				Data int `json:"data"`
			} `json:"message"`
		}
		fmt.Println("reading message")
		err := c.ReadMessage(ctx, chName, &message)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Println("Message Content -> ", message.Message.Data)

		err = c.CommitMessage(ctx, chName)
		if err != nil {
			log.Println(err.Error())
		}

	}
}
