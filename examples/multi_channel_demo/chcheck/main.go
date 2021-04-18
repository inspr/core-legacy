package main

import (
	"fmt"

	dappclient "inspr.dev/inspr/pkg/client"
	"inspr.dev/inspr/pkg/sidecar/models"
	"golang.org/x/net/context"
)

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testChannels := []string{"testch1", "testch2", "testch3"}
	checkChannel := "checkch"

	for {
		for i := 0; i < 3; i++ {
			var testMsg models.BrokerData
			err := client.ReadMessage(ctx, testChannels[i], &testMsg)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("Test received: ")
			fmt.Println(testMsg.Message.Data)

			if err := client.CommitMessage(ctx, testChannels[i]); err != nil {
				fmt.Println(err.Error())
			}
			checkMsg := models.Message{
				Data: fmt.Sprintf("%s Check!", testChannels[i]),
			}
			if err := client.WriteMessage(ctx, checkChannel, checkMsg); err != nil {
				fmt.Println(err)
				continue
			}
		}

	}
}
