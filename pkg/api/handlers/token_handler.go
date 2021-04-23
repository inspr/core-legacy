package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/inspr/inspr/pkg/auth/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
)

// TokenHandler handles requests for token creation on Insprd
func (h *Handler) TokenHandler() rest.Handler {

	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		var load models.Payload
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
		rest.JSON(w, http.StatusOK, models.JwtDO{
			Token: token,
		})

	}).Recover().Post().JSON()
}

// TokenHandler handles requests for token creation on Insprd
func (h *Handler) InitHandler() rest.Handler {

	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		load := models.Payload{
			RefreshURL: os.Getenv("REFRESH_URL"),
			Scope:      []string{""},
			Role:       1,
		}
		token, err := h.Auth.Init(load)
		if err != nil {

			rest.ERROR(w, ierrors.NewError().InternalServer().Message("unable to authenticate token").Build())
			return
		}
		rest.JSON(w, http.StatusOK, models.JwtDO{
			Token: token,
		})
	}).Recover().Post().JSON()
}