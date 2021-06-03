package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	dappclient "github.com/inspr/inspr/pkg/client"
	"github.com/inspr/inspr/pkg/sidecars/models"
)

func main() {
	counter := 5
	i := 0
	// sets up client for sidecar
	c := dappclient.NewAppClient()

	ctx, cancel := context.WithCancel(context.Background())

	message := 1234

	c.HandleChannel("mbch1", func(ctx context.Context, body io.Reader) error {
		decoder := json.NewDecoder(body)
		var testMsg models.BrokerMessage
		err := decoder.Decode(&testMsg)
		if err != nil {
			return err
		}
		i++
		fmt.Println(testMsg)
		if i >= counter {
			cancel()
			return nil
		}
		time.Sleep(time.Second * 10)
		err = c.WriteMessage(context.Background(), "mbch1", message)
		if err != nil {
			fmt.Printf("an error occurred: %v", err)
			return err
		}
		return nil
	})

	go c.Run(ctx)

	err := c.WriteMessage(context.Background(), "mbch1", message)
	if err != nil {
		fmt.Printf("an error occurred: %v", err)
		return
	}

	<-ctx.Done()
	fmt.Println("Done w/ exec")
	fmt.Scan()
}
