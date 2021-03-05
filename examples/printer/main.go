package main

import (
	"context"
	"fmt"
	"log"
	"time"

	dappclient "gitlab.inspr.dev/inspr/core/pkg/client"
)

func main() {
	// sets up client for sidecar
	c := dappclient.NewAppClient()

	// sets up ticker
	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-ticker.C:
			message, err := c.ReadMessage(context.Background(), "ch1")
			if err != nil {
				log.Println(err.Error())
			}

			fmt.Println("Message Content ", message.Data)
		}
	}
}
