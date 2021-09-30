package main

import (
	"context"
	"fmt"
	"net/http"

	"inspr.dev/inspr/examples/route_demo/model"
	dappclient "inspr.dev/inspr/pkg/client"
)

// const tCount int = 1

func main() {
	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var resp model.Response

	req := model.Request{
		Op1: 1,
		Op2: 2,
	}

	err := client.SendRequest(ctx, "api", "add", http.MethodPost, req, &resp)
	if err != nil {
		fmt.Println("ERROR!!!!!!!")
		fmt.Println(err)
		return
	}
	fmt.Println("-|-|-|-|-|-|-|-|-|-|-|-|-|-")
	fmt.Println(resp)
	// for i := 0; i < tCount; i++ {
	// 	go func() {
	// 		for {
	// 		}
	// 	}()
	// }
}
