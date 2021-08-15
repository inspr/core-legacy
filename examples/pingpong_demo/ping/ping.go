package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"context"

	dappclient "inspr.dev/inspr/pkg/client"
	"inspr.dev/inspr/pkg/sidecars/models"
)

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sentMsg := "Ping!"
	if err := client.WriteMessage(ctx, "pingoutput", sentMsg); err != nil {
		fmt.Printf("an error occurred: %v", err)
		return
	}
	client.HandleChannel(
		"pinginput",
		func(ctx context.Context, body io.Reader) error {
			var ret models.BrokerMessage

			decoder := json.NewDecoder(body)
			if err := decoder.Decode(&ret); err != nil {
				return err
			}

			fmt.Println(ret)

			if err := client.WriteMessage(ctx, "pingoutput", sentMsg); err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		},
	)
	log.Fatal(client.Run(ctx))
}
