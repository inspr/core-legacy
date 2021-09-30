package main

import (
	"encoding/json"
	"log"
	"net/http"

	"context"

	"inspr.dev/inspr/examples/route_demo/model"
	dappclient "inspr.dev/inspr/pkg/client"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/rest"
)

func addHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var data model.Request
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			rest.ERROR(w, ierrors.New("Wrong request body type").BadRequest())
			return
		}
		resp := model.Response{
			Result: data.Op1 + data.Op2,
		}
		rest.JSON(w, http.StatusOK, resp)
	}
}

func subHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var data model.Request
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			rest.ERROR(w, ierrors.New("Wrong request body type").BadRequest())
			return
		}
		resp := model.Response{
			Result: data.Op1 - data.Op2,
		}
		rest.JSON(w, http.StatusOK, resp)
	}
}

func mulHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		var data model.Request
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			rest.ERROR(w, ierrors.New("Wrong request body type").BadRequest())
			return
		}
		resp := model.Response{
			Result: data.Op1 * data.Op2,
		}
		rest.JSON(w, http.StatusOK, resp)
	}
}

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client.HandleRoute("add", addHandler())
	client.HandleRoute("sub", subHandler())
	client.HandleRoute("mul", mulHandler())
	log.Fatal(client.Run(ctx))
}
