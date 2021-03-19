package main

import (
	"fmt"

	dappclient "gitlab.inspr.dev/inspr/core/pkg/client"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
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
			testMsg := models.Message{
				Data: fmt.Sprintf("Testing channel: %s", testChannels[i]),
			}
			if err := client.WriteMessage(ctx, testChannels[i], testMsg); err != nil {
				fmt.Println(err)
				continue
			}
			var checkMsg models.Message
			err := client.ReadMessage(ctx, checkChannel, &checkMsg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Check received: ")
			fmt.Println(checkMsg.Data)
			if err := client.CommitMessage(ctx, checkChannel); err != nil {
				fmt.Println(err.Error())
			}
		}

	}
}
