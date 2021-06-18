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

	sentMsg := "Pong!"
	if err := client.WriteMessage(ctx, "pongoutput", sentMsg); err != nil {
		fmt.Printf("an error occurred: %v", err)
		return
	}
	client.HandleChannel("ponginput", func(ctx context.Context, body io.Reader) error {
		var ret models.BrokerMessage

		decoder := json.NewDecoder(body)
		if err := decoder.Decode(&ret); err != nil {
			return err
		}

		fmt.Println(ret)

		if err := client.WriteMessage(ctx, "pongoutput", sentMsg); err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	log.Fatal(client.Run(ctx))
}
