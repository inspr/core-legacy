package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"context"

	dappclient "inspr.dev/inspr/pkg/client"
	"inspr.dev/inspr/pkg/sidecars/models"
)

// func addHandler() rest.Handler {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var data model.Request
// 		err := json.NewDecoder(r.Body).Decode(&data)
// 		if err != nil {
// 			rest.ERROR(w, ierrors.New("Wrong request body type").BadRequest())
// 			return
// 		}
// 		resp := model.Response{
// 			Result: data.Op1 + data.Op2,
// 		}
// 		rest.JSON(w, http.StatusOK, resp)
// 	}
// }

// func subHandler() rest.Handler {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var data model.Request
// 		err := json.NewDecoder(r.Body).Decode(&data)
// 		if err != nil {
// 			rest.ERROR(w, ierrors.New("Wrong request body type").BadRequest())
// 			return
// 		}
// 		resp := model.Response{
// 			Result: data.Op1 - data.Op2,
// 		}
// 		rest.JSON(w, http.StatusOK, resp)
// 	}
// }

// func mulHandler() rest.Handler {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var data model.Request
// 		err := json.NewDecoder(r.Body).Decode(&data)
// 		if err != nil {
// 			rest.ERROR(w, ierrors.New("Wrong request body type").BadRequest())
// 			return
// 		}
// 		resp := model.Response{
// 			Result: data.Op1 * data.Op2,
// 		}
// 		rest.JSON(w, http.StatusOK, resp)
// 	}
// }

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client.HandleChannel("receivech", func(ctx context.Context, body io.Reader) error {
		var ret models.BrokerMessage

		decoder := json.NewDecoder(body)
		if err := decoder.Decode(&ret); err != nil {
			return err
		}

		fmt.Println(ret)
		return nil
	})
	log.Fatal(client.Run(ctx))
}
