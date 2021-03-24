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

var discordCH = "discodMessages"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Mux.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		data := Message{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		if data.Discord {
			discordMsg := models.Message{
				Data: data.Message,
			}
			if err := client.WriteMessage(ctx, discordCH, discordMsg); err != nil {
				fmt.Println(err)
				rest.ERROR(w, err)
			}
		}

		rest.JSON(w, http.StatusOK, nil)
	})
}

// Run starts the server on the port given in addr
func (s *Server) Run(addr string) {
	fmt.Printf("MultiMessage api is up! Listening on port: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.Mux))
}
