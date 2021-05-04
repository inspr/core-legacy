package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	dappclient "github.com/inspr/inspr/pkg/client"
	"golang.org/x/net/context"
)

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sentMsg := "ping"
	if err := client.WriteMessage(ctx, "pingoutput", sentMsg); err != nil {
		fmt.Println(err)
		return
	}
	client.HandleChannel("pinginput", func(_ context.Context, r io.Reader) error {
		decoder := json.NewDecoder(r)
		var ret struct{ Message string }
		err := decoder.Decode(&ret.Message)

		if err != nil {
			return err
		}

		if err := client.WriteMessage(ctx, "pingoutput", sentMsg); err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	log.Fatal(client.Run(ctx))
}
