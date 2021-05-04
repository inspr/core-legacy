package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"

	dappclient "github.com/inspr/inspr/pkg/client"
)

func main() {

	// sets up client for sidecar
	client := dappclient.NewAppClient()

	// channelName
	inputChannel := "input"
	outputChannel := "output"

	type Message struct {
		Message int    `json:"message"`
		Channel string `json:"channel"`
	}

	fmt.Println("starting...")
	// handles messages sent to the input channel
	client.HandleChannel(inputChannel, func(_ context.Context, r io.Reader) error {
		decoder := json.NewDecoder(r)
		var msg Message

		err := decoder.Decode(&msg)
		if err != nil {
			return err
		}
		log.Printf("msg.Message = %+v\n", msg.Message)
		if big.NewInt(int64(msg.Message)).ProbablyPrime(0) {
			err = client.WriteMessage(
				context.Background(),
				outputChannel,
				msg.Message,
			)
			if err != nil {
				return err
			}
		}
		// writes a message in the output channel
		return nil
	})
	log.Fatalln(client.Run(context.Background()))
}
