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
	testChannels := []string{"testch1", "testch2", "testch3"}
	checkChannel := "checkch"

	client.HandleChannel(
		checkChannel,
		func(ctx context.Context, body io.Reader) error {
			decoder := json.NewDecoder(body)
			var checkMsg models.BrokerMessage
			err := decoder.Decode(&checkMsg)
			if err != nil {
				return err
			}

			log.Println("Check received")
			log.Println(checkMsg.Data)
			return nil
		},
	)

	go func() {
		log.Fatalln(client.Run(ctx))
	}()
	for {
		for _, testChannel := range testChannels {
			testMsg := fmt.Sprintf("Testing channel: %s", testChannel)
			if err := client.WriteMessage(ctx, testChannel, testMsg); err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
