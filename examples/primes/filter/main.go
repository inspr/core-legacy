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
	var number float64

	// sets up ticker to sync with generator
	ticker := time.NewTicker(2 * time.Second)
	rand.Seed(time.Now().UnixNano())

	// sets up client for sidecar
	client := dappclient.NewAppClient()

	// channelName
	inputChannel := "primes_ch1"
	outputChannel := "primes_ch2"

	for {
		select {
		case <-ticker.C:
			msg, err := client.ReadMessage(context.Background(), inputChannel)
			if err != nil {
				fmt.Println(err.Error())
			}
			number = msg.Data.(float64)
			fmt.Println("Read: ", number)

			err = client.CommitMessage(context.Background(), inputChannel)
			if err != nil {
				fmt.Println(err.Error())
			}

			if big.NewInt(int64(number)).ProbablyPrime(0) {
				client.WriteMessage(context.Background(), outputChannel, models.Message{
					Data: number,
				})
			}
		}
	}
}
