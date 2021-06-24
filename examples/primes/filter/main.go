package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"

	dappclient "inspr.dev/inspr/pkg/client"
	"inspr.dev/inspr/pkg/sidecars/models"
)

func main() {

	// sets up client for sidecar
	client := dappclient.NewAppClient()

	// channelName
	inputChannel := "input"
	outputChannel := "output"

	fmt.Println("starting...")
	// handles messages sent to the input channel
	client.HandleChannel(inputChannel, func(_ context.Context, r io.Reader) error {
		decoder := json.NewDecoder(r)
		var msg models.BrokerMessage

		err := decoder.Decode(&msg)
		if err != nil {
			return err
		}

		log.Printf("Message = %+v\n", msg.Data)
		msgNumber, ok := msg.Data.(int64)
		if !ok {
			return fmt.Errorf("unable to convert '%v' to int64", msg.Data)
		}

		if big.NewInt(msgNumber).ProbablyPrime(0) {
			if err := client.WriteMessage(context.Background(), outputChannel, msg.Data); err != nil {
				return err
			}
		}
		// writes a message in the output channel
		return nil
	})
	log.Fatalln(client.Run(context.Background()))
}
