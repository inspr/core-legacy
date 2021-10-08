package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"inspr.dev/inspr/examples/route_demo/model"
	dappclient "inspr.dev/inspr/pkg/client"
)

func main() {
	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		client.Run(ctx)
	}()

	var resp model.Response
	var err error

	req := model.Request{
		Op1: 1,
		Op2: 2,
	}

	time.Sleep(5 * time.Second)
	for {
		err = client.SendRequest(ctx, "api", "add", http.MethodPost, req, &resp)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp.Result)
		}
		req.Op1 = resp.Result
		err = client.SendRequest(ctx, "api", "mul", http.MethodPost, req, &resp)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp.Result)
		}
		req.Op1 = resp.Result
		err = client.SendRequest(ctx, "api", "sub", http.MethodPost, req, &resp)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp.Result)
		}
	}
}
