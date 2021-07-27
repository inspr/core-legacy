package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/rest"
)

var initialized bool

// HandleInit handles initialization on the server
func (server *Server) HandleInit() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Key string
			auth.Payload
		}

		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&data)
		if initialized {
			rest.ERROR(w, ierrors.New("already initialized").AlreadyExists())
			return
		}

		server.logger.Debug("received data to initialize auth service",
			zap.Any("data: ", data))

		if data.Key != os.Getenv("INSPR_INIT_KEY") {
			rest.ERROR(w, ierrors.New("invalid key").Forbidden())
			return
		}

		initialized = true
		token, err := server.tokenize(
			data.Payload,
			time.Now().Add(time.Minute*30),
		)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		rest.JSON(w, 200, auth.JwtDO{Token: token})

	}).Post().JSON().Recover()
}
