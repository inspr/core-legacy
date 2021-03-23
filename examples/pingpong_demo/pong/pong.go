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

	for {
		sentMsg := models.Message{
			Data: "Pong!",
		}

		if err := client.WriteMessage(ctx, "ppChannel2", sentMsg); err != nil {
			fmt.Println(err)
			continue
		}

		var recMsg models.Message
		err := client.ReadMessage(ctx, "ppChannel1", &recMsg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Read message: ")
		fmt.Println(recMsg.Data)

		if err := client.CommitMessage(ctx, "ppChannel1"); err != nil {
			fmt.Println(err.Error())
		}
	}
}
