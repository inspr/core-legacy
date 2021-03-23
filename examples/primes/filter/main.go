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
	var number int

	// sets up ticker to sync with generator
	ticker := time.NewTicker(2 * time.Second)
	rand.Seed(time.Now().UnixNano())

	// sets up client for sidecar
	client := dappclient.NewAppClient()

	// channelName
	inputChannel := "input"
	outputChannel := "output"

	type Message struct {
		Message struct {
			Data int `json:"data"`
		} `json:"message"`
		Channel string `json:"channel"`
	}
	fmt.Println("starting...")
	for range ticker.C {
		var msg Message
		fmt.Println("reading message...")
		err := client.ReadMessage(context.Background(), inputChannel, &msg)
		if err != nil {
			fmt.Println(err.Error())
		}

		number = msg.Message.Data
		fmt.Println("Read: ", number)

		err = client.CommitMessage(context.Background(), inputChannel)
		if err != nil {
			fmt.Println(err.Error())
		}

		if big.NewInt(int64(number)).ProbablyPrime(0) {
			client.WriteMessage(
				context.Background(),
				outputChannel,
				models.Message{
					Data: number,
				},
			)
		}
	}
}
