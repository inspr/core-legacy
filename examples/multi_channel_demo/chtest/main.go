package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	dappclient "github.com/inspr/inspr/pkg/client"
	"github.com/inspr/inspr/pkg/sidecar/models"
	"golang.org/x/net/context"
)

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testChannels := []string{"testch1", "testch2", "testch3"}
	checkChannel := "checkch"

	for i := 0; i < 3; i++ {
		client.HandleChannel(checkChannel, func(ctx context.Context, body io.Reader) error {
			decoder := json.NewDecoder(body)
			var checkMsg models.BrokerData
			err := decoder.Decode(&checkMsg)
			if err != nil {
				return err
			}

			log.Println("Check received")
			log.Println(checkMsg.Message)
			return nil
		})

	}
	go func() {
		log.Fatalln(client.Run(ctx))
	}()
	for {

		for i := 0; i < 3; i++ {
			testMsg := fmt.Sprintf("Testing channel: %s", testChannels[i])
			if err := client.WriteMessage(ctx, testChannels[i], testMsg); err != nil {
				fmt.Println(err)
				continue
			}
		}

	}
}
