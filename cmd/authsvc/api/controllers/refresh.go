package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/auth/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

// Refresh returns the refreshing endpoint. This entpoint receives a refresh token and a refresh url, it returns a refreshed token.
func (server *Server) Refresh() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContent := r.Header["Authorization"]

		if len(headerContent) != 1 ||
			!strings.HasPrefix(headerContent[0], "Bearer ") {
			http.Error(
				w,
				"Bad Request, expected: Authorization: Bearer <token>",
				http.StatusUnauthorized,
			)
			return
		}

		token := []byte(strings.TrimPrefix(headerContent[0], "Bearer "))

		_, err := jwt.Parse(
			[]byte(token),
			jwt.WithValidate(true),
			jwt.WithVerify(jwa.RS256, server.privKey.PublicKey),
		)
		if err != nil {
			http.Error(
				w,
				"Invalid token",
				http.StatusForbidden,
			)
			return
		}

		load, err := auth.Desserialize(token)
		if err != nil {
			http.Error(
				w,
				"Invalid token",
				http.StatusForbidden,
			)
			return
		}

		payload, err := refreshPayload(load.Refresh, load.RefreshURL)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		signed, err := server.tokenize(*payload)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		respBody := models.JwtDO{
			Token: signed,
		}
		rest.JSON(w, http.StatusOK, respBody)
	}
}

func refreshPayload(refresToken []byte, refreshURL string) (*models.Payload, error) {
	reqBody := models.ResfreshDO{
		RefreshToken: refresToken,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		err = ierrors.NewError().InternalServer().Message(err.Error()).Build()
		return nil, err
	}

	c := &http.Client{}
	resp, err := c.Post(refreshURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil || resp.StatusCode != http.StatusOK {
		err = ierrors.NewError().InternalServer().InnerError(err).Build()
		return nil, err
	}
	defer resp.Body.Close()

	payload := models.Payload{}
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
