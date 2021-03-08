package main

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	dappclient "gitlab.inspr.dev/inspr/core/pkg/client"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

func main() {
	var number int64

	// sets up ticker to sync with generator
	ticker := time.NewTicker(2 * time.Second)
	rand.Seed(time.Now().UnixNano())

	// sets up client for sidecar
	c := dappclient.NewAppClient()

	// channelName
	inputChannel := "ch1"
	outputChannel := "ch2"

	for {
		select {
		case <-ticker.C:
			msg, err := c.ReadMessage(context.Background(), inputChannel)
			if err != nil {
				fmt.Println(err.Error())
			}
			number = msg.Data.(int64)
			c.CommitMessage(context.Background(), inputChannel)

			if big.NewInt(number).ProbablyPrime(0) {
				c.WriteMessage(context.Background(), outputChannel, models.Message{
					Data: number,
				})
			}
		}
	}
}
