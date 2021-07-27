package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/rest"
)

// Tokenize receives a token's payload and encodes it in a jwt
func (server *Server) Tokenize() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		data := auth.Payload{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			server.logger.Error("unable to decode payload",
				zap.Any("error", err))

			err = ierrors.Wrap(
				ierrors.From(err).BadRequest(),
				"invalid body",
			)

			rest.ERROR(w, err)
			return
		}

		signed, err := server.tokenize(data, time.Now().Add(time.Minute*1))
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		body := auth.JwtDO{
			Token: signed,
		}
		rest.JSON(w, http.StatusOK, body)
	}
}

func (server *Server) tokenize(payload auth.Payload, exp time.Time) ([]byte, error) {
	var err error
	token := jwt.New()
	token.Set(jwt.ExpirationKey, exp)
	token.Set("payload", payload)

	signed, err := jwt.Sign(token, jwa.RS256, server.privKey)
	if err != nil {
		server.logger.Error("unable to sign JWT with provided RSA private key", zap.Any("error", err))
		err := ierrors.New("unable to sign JWT with available RSA private key").InternalServer()
		return nil, err
	}
	return signed, nil
}
