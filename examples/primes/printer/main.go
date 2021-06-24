package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	dappclient "inspr.dev/inspr/pkg/client"
	"inspr.dev/inspr/pkg/sidecars/models"
)

func main() {
	// sets up client for sidecar
	c := dappclient.NewAppClient()

	// sets up ticker
	chName := "input"
	fmt.Println("starting...")
	c.HandleChannel(chName, func(_ context.Context, r io.Reader) error {

		var message models.BrokerMessage

		fmt.Println("reading message")
		decoder := json.NewDecoder(r)
		decoder.Decode(&message)
		fmt.Println("the number ", message.Data, " is a prime")
		return nil

	})

	log.Fatalln(c.Run(context.Background()))
}
