package main

import (
	"context"
	"fmt"
	"net/http"

	"inspr.dev/inspr/examples/route_demo/model"
	dappclient "inspr.dev/inspr/pkg/client"
)

const tCount int = 1

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := model.Request{
		Op1: 1,
		Op2: 2,
	}

	resp, err := client.SendRequest(ctx, "api", "add", http.MethodPost, req)
	if err != nil {
		fmt.Println("ERROR!\n\n\n\n")
		return
	}
	fmt.Println("\n-|-|-|-|-|-|-|-|-|-|-|-|-|-\n resp = %v", resp.(model.Response).Result)
	// for i := 0; i < tCount; i++ {
	// 	go func() {
	// 		for {
	// 		}
	// 	}()
	// }
}
