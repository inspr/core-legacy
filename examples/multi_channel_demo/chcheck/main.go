package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	dappclient "github.com/inspr/inspr/pkg/client"
	"github.com/inspr/inspr/pkg/sidecars/models"
	"golang.org/x/net/context"
)

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testChannels := []string{"testch1", "testch2", "testch3"}
	checkChannel := "checkch"

	for i := 0; i < 3; i++ {
		testChannel := testChannels[i]
		client.HandleChannel(testChannel, func(ctx context.Context, body io.Reader) error {
			decoder := json.NewDecoder(body)
			var testMsg models.BrokerMessage
			err := decoder.Decode(&testMsg)
			if err != nil {
				return err
			}

			checkMessage := fmt.Sprintf("%s Check!", testChannel)
			if err := client.WriteMessage(ctx, checkChannel, checkMessage); err != nil {
				return err
			}
			return nil
		})

	}
	log.Fatalln(client.Run(ctx))
}
