package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"strconv"

	dappclient "inspr.dev/inspr/pkg/client"
	"inspr.dev/inspr/pkg/sidecars/models"
)

func main() {

	// sets up client for sidecar
	client := dappclient.NewAppClient()
	ctx := context.Background()

	// channelName
	inputChannel := "filterinput"
	outputChannel := "filteroutput"

	fmt.Println("starting...")
	// handles messages sent to the input channel
	client.HandleChannel(inputChannel, func(ctx context.Context, r io.Reader) error {
		var msg models.BrokerMessage
		decoder := json.NewDecoder(r)

		err := decoder.Decode(&msg)
		if err != nil {
			return err
		}

		fmt.Printf("Message: %+v\n", msg.Data)
		strMsg := fmt.Sprintf("%v", msg.Data)
		msgNumber, err := strconv.ParseInt(strMsg, 10, 64)
		if err != nil {
			fmt.Printf("unable to convert '%v' to int64: %v", msg.Data, err)
			return err
		}

		if big.NewInt(msgNumber).ProbablyPrime(0) {
			if err := client.WriteMessage(context.Background(), outputChannel, msg.Data); err != nil {
				fmt.Println(err)
				return err
			}
		}
		// writes a message in the output channel
		return nil
	})
	log.Fatalln(client.Run(ctx))
}
