package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/auth/models"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// Refresh returns the refreshing endpoint. This entpoint receives a refresh token and a refresh url, it returns a refreshed token.
func (server *Server) Refresh() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		data := models.ResfreshDI{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			err = ierrors.NewError().BadRequest().Message("invalid body").Build()
			rest.ERROR(w, err)
			return
		}

		reqBody := models.ResfreshDO{
			RefreshToken: data.RefreshToken,
		}
		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			err = ierrors.NewError().BadRequest().Message("invalid request body").Build()
			rest.ERROR(w, err)
			return
		}

		c := &http.Client{}
		resp, err := c.Post(data.RefreshURL, "application/json", bytes.NewBuffer(reqBytes))
		if err != nil || resp.StatusCode != http.StatusOK {
			rest.ERROR(w, err)
			return
		}

		payload := models.Payload{}
		err = json.NewDecoder(resp.Body).Decode(&payload)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		signed, err := server.tokenize(payload)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		body := models.JwtDO{
			Token: signed,
		}
		rest.JSON(w, http.StatusOK, body)
	}
}
