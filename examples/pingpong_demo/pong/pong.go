package main

import (
	"fmt"

	dappclient "gitlab.inspr.dev/inspr/core/pkg/client"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
	"golang.org/x/net/context"
)

type expectedDataType struct {
	Message struct {
		Data string `json:"data"`
	} `json:"message"`
	Channel string `json:"channel"`
}

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		sentMsg := models.Message{
			Data: "Pong!",
		}

		if err := client.WriteMessage(ctx, "pongoutput", sentMsg); err != nil {
			fmt.Println(err)
			continue
		}

		var recMsg expectedDataType
		err := client.ReadMessage(ctx, "ponginput", &recMsg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Read message: ")
		fmt.Println(recMsg.Message.Data)

		if err := client.CommitMessage(ctx, "ponginput"); err != nil {
			fmt.Println(err.Error())
		}
	}
}
