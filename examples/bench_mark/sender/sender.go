package main

import (
	"fmt"
	"time"

	"context"

	dappclient "inspr.dev/inspr/pkg/client"
)

const tCount int = 1

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sentMsg := "Ping!"

	errCh := make(chan error)

	for i := 0; i < tCount; i++ {
		go func() {
			for {
				if err := client.WriteMessage(ctx, "sendch", sentMsg); err != nil {
					fmt.Printf("an error occurred: %v", err)
					errCh <- err
					return
				}
				select {
				case err := <-errCh:
					fmt.Printf("an error occurred: %v", err)
					return
				default:
				}
			}
		}()
	}

	<-errCh
	<-time.After(5 * time.Second)
}
