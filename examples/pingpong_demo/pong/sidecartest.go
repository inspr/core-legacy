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
			Data: "Boa tarde amigos",
		}

		if err := client.WriteMessage(ctx, "ch2", sentMsg); err != nil {
			fmt.Println(err)
			continue
		}

		recMsg, err := client.ReadMessage(ctx, "ch1")
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Mensagem lida: ")
		fmt.Println(recMsg.Data)

		if err := client.CommitMessage(ctx, "ch1"); err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}
