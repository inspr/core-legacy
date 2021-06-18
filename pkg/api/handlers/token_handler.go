package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/rest"
)

// ControllerRefreshHandler handles requests for token refresing on inspr controllers on Insprd
func (h *Handler) ControllerRefreshHandler() rest.Handler {

	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		received := auth.ResfreshDO{}
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			log.Printf("err = %+v\n", err)
			rest.ERROR(w, err)
			return
		}

		// this is the path to the app
		appQuery := string(received.RefreshToken)

		app, err := h.Memory.Root().Apps().Get(appQuery)
		if err != nil {
			log.Printf("err = %+v\n", err)
			rest.ERROR(w, err)
			return
		}

		// refresh the payload with the current permissions of the dApp
		payload := auth.Payload{
			UID: app.Meta.UUID,
			Permissions: map[string][]string{
				app.Spec.Auth.Scope: app.Spec.Auth.Permissions,
			},
			Refresh:    []byte(appQuery),
			RefreshURL: fmt.Sprintf("%v/refreshController", os.Getenv("INSPR_INSPRD_ADDRESS")),
		}
		rest.JSON(w, 200, payload)

	}).Recover().Post().JSON()
}

// TokenHandler handles requests for token creation on Insprd
func (h *Handler) TokenHandler() rest.Handler {

	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		var load auth.Payload
		err := json.NewDecoder(r.Body).Decode(&load)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		token, err := h.Auth.Tokenize(load)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		rest.JSON(w, http.StatusOK, auth.JwtDO{
			Token: token,
		})

	}).Recover().Post().JSON()
}

// InitHandler handles requests for cluster auth initialization
func (h *Handler) InitHandler() rest.Handler {

	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		res := struct{ Key string }{}
		decoder.Decode(&res)
		load := auth.Payload{
			RefreshURL:  os.Getenv("REFRESH_URL"),
			Permissions: map[string][]string{"": {auth.CreateToken}},
		}
		token, err := h.Auth.Init(res.Key, load)
		if err != nil {
			rest.ERROR(w, ierrors.NewError().InternalServer().Message("unable to authenticate token: %v", err.Error()).Build())
			return
		}
		rest.JSON(w, http.StatusOK, auth.JwtDO{
			Token: token,
		})
	}).Recover().Post().JSON()
}
