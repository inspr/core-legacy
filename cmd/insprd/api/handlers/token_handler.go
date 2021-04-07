package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/auth"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

func (h *Handler) TokenHandler() rest.Handler {
	type TokenReturn struct {
		Token string `json:"token"`
	}
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		var load auth.Payload
		err := json.NewDecoder(r.Body).Decode(&load)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		token, err := h.auth.Tokenize(load)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		rest.JSON(w, http.StatusOK, TokenReturn{
			token,
		})

	}).Post().JSON()
}
