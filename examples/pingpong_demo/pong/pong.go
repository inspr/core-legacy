package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"golang.org/x/net/context"
	dappclient "inspr.dev/inspr/pkg/client"
)

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sentMsg := "pong"
	if err := client.WriteMessage(ctx, "pongoutput", sentMsg); err != nil {
		fmt.Println(err)
		return
	}
	client.HandleChannel("ponginput", func(_ context.Context, r io.Reader) error {
		decoder := json.NewDecoder(r)
		var ret struct{ Message string }
		err := decoder.Decode(&ret.Message)

		if err != nil {
			return err
		}

		if err := client.WriteMessage(ctx, "pongoutput", sentMsg); err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	log.Fatal(client.Run(ctx))
}
