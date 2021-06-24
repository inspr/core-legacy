package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	dappclient "inspr.dev/inspr/pkg/client"
	"inspr.dev/inspr/pkg/sidecars/models"
)

func main() {
	maxIterations := 500
	counter := 0
	// sets up client for sidecar
	c := dappclient.NewAppClient()

	ctx, cancel := context.WithCancel(context.Background())

	c.HandleChannel("mbch1", func(ctx context.Context, body io.Reader) error {
		decoder := json.NewDecoder(body)
		var testMsg models.BrokerMessage
		err := decoder.Decode(&testMsg)
		if err != nil {
			return err
		}
		counter++
		fmt.Println(testMsg)
		if counter >= maxIterations {
			cancel()
			return nil
		}
		time.Sleep(time.Second * 5)
		err = c.WriteMessage(context.Background(), "mbch1", counter)
		if err != nil {
			fmt.Printf("an error occurred: %v", err)
			return err
		}
		return nil
	})

	go c.Run(ctx)

	err := c.WriteMessage(context.Background(), "mbch1", counter)
	if err != nil {
		fmt.Printf("an error occurred: %v", err)
		return
	}

	<-ctx.Done()
	fmt.Println("Done w/ exec")
}
