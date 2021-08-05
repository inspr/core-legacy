package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/rest"
)

// Refresh returns the refreshing endpoint. This entpoint receives a refresh token and a refresh url, it returns a refreshed token.
func (server *Server) Refresh() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContent := r.Header["Authorization"]

		if len(headerContent) != 1 ||
			!strings.HasPrefix(headerContent[0], "Bearer ") {
			err := ierrors.New(
				"bad Request, expected: Authorization: Bearer <token>",
			).Unauthorized()
			rest.ERROR(w, err)
			return
		}

		token := []byte(strings.TrimPrefix(headerContent[0], "Bearer "))

		server.logger.Info("parsing received bearer token")
		_, err := jwt.Parse(
			token,
			jwt.WithValidate(true),
			jwt.WithVerify(jwa.RS256, server.privKey.PublicKey),
		)
		if err != nil && err.Error() != `exp not satisfied` {
			err := ierrors.Wrap(
				ierrors.New(err).Forbidden(),
				"couldn't parse token",
			)
			rest.ERROR(w, err)
			return
		}

		server.logger.Info("deserializing parsed token")
		load, err := auth.Desserialize(token)
		if err != nil {
			err := ierrors.Wrap(
				err,
				"couldn't desserialize token",
			)
			rest.ERROR(w, err)
			return
		}

		server.logger.Debug("received payload", zap.Any("content", load))

		server.logger.Info("refreshing old payload")
		payload, err := refreshPayload(load.Refresh, load.RefreshURL)
		if err != nil {
			err := ierrors.Wrap(
				err,
				"couldn't refresh payload",
			)
			rest.ERROR(w, err)
			return
		}

		server.logger.Debug("refreshed payload", zap.Any("content", payload))

		signed, err := server.tokenize(*payload, time.Now().Add(time.Minute*8))
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		server.logger.Debug("new token", zap.String("value", string(signed)))

		respBody := auth.JwtDO{
			Token: signed,
		}

		rest.JSON(w, http.StatusOK, respBody)
	}
}

func refreshPayload(refreshToken []byte, refreshURL string) (*auth.Payload, error) {
	reqBody := auth.ResfreshDO{
		RefreshToken: refreshToken,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		err = ierrors.New(err).InternalServer()
		return nil, err
	}

	resp, err := http.Post(refreshURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil || resp.StatusCode != http.StatusOK {
		err = ierrors.New(err).InternalServer()
		return nil, err
	}
	defer resp.Body.Close()

	payload := auth.Payload{}
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
