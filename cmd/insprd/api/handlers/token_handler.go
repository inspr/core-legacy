package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/auth/models"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
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
		token, err := h.auth.Tokenize(load)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		rest.JSON(w, http.StatusOK, models.JwtDO{
			Token: token,
		})

	}).Post().JSON()
}
