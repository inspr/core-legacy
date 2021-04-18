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

	sentMsg := models.Message{
		Data: "Boa tarde amigos",
	}

	if err := client.WriteMessage(ctx, "ch1", sentMsg); err != nil {
		fmt.Println(err)
		return
	}

	var recMsg models.Message
	err := client.ReadMessage(ctx, "ch1", &recMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Mensagem lida: ")
	fmt.Println(recMsg.Data)

	if err := client.CommitMessage(ctx, "ch1"); err != nil {
		fmt.Println(err.Error())
		return
	}
}
