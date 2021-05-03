package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

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

		log.Printf("headerContent = %+v\n", headerContent)
		if len(headerContent) != 1 ||
			!strings.HasPrefix(headerContent[0], "Bearer ") {
			err := ierrors.NewError().Unauthorized().Message("bad Request, expected: Authorization: Bearer <token>").Build()
			rest.ERROR(w, err)
			return
		}

		token := []byte(strings.TrimPrefix(headerContent[0], "Bearer "))
		log.Printf("string(token) = %+v\n", string(token))

		_, err := jwt.Parse(
			token,
			jwt.WithValidate(true),
			jwt.WithVerify(jwa.RS256, server.privKey.PublicKey),
		)
		if err != nil && err.Error() != `exp not satisfied` {
			err := ierrors.NewError().Forbidden().Message("invalid token").Build()
			rest.ERROR(w, err)
			return
		}

		load, err := auth.Desserialize(token)
		if err != nil {
			err := ierrors.NewError().Forbidden().Message("invalid token, error: %s", err.Error()).Build()
			rest.ERROR(w, err)
			return
		}
		log.Printf("load = %+v\n", load)

		payload, err := refreshPayload(load.Refresh, load.RefreshURL)
		if err != nil {
			err := ierrors.NewError().InternalServer().Message("invalid token").Build()
			rest.ERROR(w, err)
			return
		}
		log.Printf("payload = %+v\n", payload)

		signed, err := server.tokenize(*payload, time.Now().Add(time.Minute*30))
		if err != nil {
			err := ierrors.NewError().InternalServer().Message(err.Error()).Build()
			rest.ERROR(w, err)
			return
		}
		log.Printf("string(signed) = %+v\n", string(signed))

		respBody := models.JwtDO{
			Token: signed,
		}

		rest.JSON(w, http.StatusOK, respBody)
	}
}

func refreshPayload(refreshToken []byte, refreshURL string) (*models.Payload, error) {
	reqBody := models.ResfreshDO{
		RefreshToken: refreshToken,
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
