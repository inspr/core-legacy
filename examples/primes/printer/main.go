package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	dappclient "github.com/inspr/inspr/pkg/client"
)

func main() {
	// sets up client for sidecar
	c := dappclient.NewAppClient()

	// sets up ticker
	chName := "input"
	fmt.Println("starting...")
	c.HandleChannel(chName, func(_ context.Context, r io.Reader) error {

		var message struct {
			Message int `json:"message"`
		}
		fmt.Println("reading message")
		decoder := json.NewDecoder(r)
		decoder.Decode(&message)
		fmt.Println("the number ", message.Message, " is a prime")
		return nil

	})

	log.Fatalln(c.Run(context.Background()))
}
