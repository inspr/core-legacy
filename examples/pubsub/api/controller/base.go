package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dappclient "gitlab.inspr.dev/inspr/core/pkg/client"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	Mux *http.ServeMux
}

var discordCH = "pubsubch"

type Message struct {
	Message string `json:"message"`
	Discord bool   `json:"discord"`
	Slack   bool   `json:"slack"`
	Twitter bool   `json:"twitter"`
}

// Init - configures the server
func (s *Server) Init() {
	s.Mux = http.NewServeMux()
	client := dappclient.NewAppClient()
	s.Mux.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		data := Message{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		discordMsg := models.Message{
			Data: data.Message,
		}
		if err := client.WriteMessage(ctx, discordCH, discordMsg); err != nil {
			fmt.Println(err)
			rest.ERROR(w, err)
		}

		rest.JSON(w, http.StatusOK, nil)
	})
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) { // this is called by the main()
	fmt.Printf("pubsub api is up! Listening on port: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.Mux))
}
