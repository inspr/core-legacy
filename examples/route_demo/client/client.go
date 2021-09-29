package main

import (
	"context"

	dappclient "inspr.dev/inspr/pkg/client"
)

const tCount int = 1

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < tCount; i++ {
		go func() {
			for {
			}
		}()
	}
}
