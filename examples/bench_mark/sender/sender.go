package main

import (
	"fmt"

	"context"

	dappclient "inspr.dev/inspr/pkg/client"
)

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sentMsg := "Ping!"

	for {
		if err := client.WriteMessage(ctx, "sendch", sentMsg); err != nil {
			fmt.Printf("an error occurred: %v", err)
			return
		}
	}

}
