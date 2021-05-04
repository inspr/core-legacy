package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
)

// Tokenize receives a token's payload and encodes it in a jwt
func (server *Server) Tokenize() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		data := auth.Payload{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			server.logger.Error("unable to decode ")
			err = ierrors.NewError().BadRequest().Message("invalid body, error: %s", err.Error()).Build()
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
		err := ierrors.NewError().InternalServer().Message("unable to sign JWT with availible RSA private key").Build()
		return nil, err
	}
	return signed, nil
}
