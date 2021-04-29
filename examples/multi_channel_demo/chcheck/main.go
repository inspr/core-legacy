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
		client.HandleChannel(testChannels[i], func(ctx context.Context, body io.Reader) error {
			decoder := json.NewDecoder(body)
			var testMsg models.BrokerData
			err := decoder.Decode(&testMsg)
			if err != nil {
				return err
			}

			checkMessage := fmt.Sprintf("%s Check!", testChannels[i])
			if err := client.WriteMessage(ctx, checkChannel, checkMessage); err != nil {
				return err
			}
			return nil
		})

	}
	go func() {
		log.Fatalln(client.Run(ctx))
	}()
}
