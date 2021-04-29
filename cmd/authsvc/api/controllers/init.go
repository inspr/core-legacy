package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/inspr/inspr/pkg/auth/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
)

var initialized bool

// HandleInit handles initialization on the server
func (server *Server) HandleInit() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Key string
			models.Payload
		}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&data)
		if initialized {
			rest.ERROR(w, ierrors.NewError().Message("already initialized").Build())
			return
		}
		log.Printf("data = %+v\n", data)
		if data.Key != os.Getenv("INSPR_INIT_KEY") {
			rest.ERROR(w, ierrors.NewError().Message("invalid key").Forbidden().Build())
			return
		}
		initialized = true
		token, err := server.tokenize(data.Payload, time.Now().Add(time.Minute*30))
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		rest.JSON(w, 200, models.JwtDO{Token: token})

	}).Post().JSON().Recover()
}
