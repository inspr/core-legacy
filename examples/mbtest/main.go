package main

import (
	"context"
	"fmt"
	"time"

	dappclient "github.com/inspr/inspr/pkg/client"
)

func main() {
	counter := 1000
	// sets up client for sidecar
	c := dappclient.NewAppClient()

	message := 1234

	for i := 0; i < counter; i++ {
		err := c.WriteMessage(context.Background(), "mbch1", message)
		if err != nil {
			fmt.Printf("an error occurred: %v", err)
		}
		time.Sleep(time.Second * 5)
	}
}
