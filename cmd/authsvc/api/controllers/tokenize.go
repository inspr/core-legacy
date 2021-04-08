package controllers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"gitlab.inspr.dev/inspr/core/cmd/authsvc/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/auth"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
	"go.uber.org/zap"
)

// Tokenize receives a token's payload and encodes it in a jwt
func (server *Server) Tokenize() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		data := auth.Payload{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			err = ierrors.NewError().BadRequest().Message("invalid body").Build()
			rest.ERROR(w, err)
			return
		}

		token := jwt.New()
		token.Set(jwt.ExpirationKey, time.Now().Add(30*time.Minute))
		token.Set("payload", data)

		keyPem := []byte(os.Getenv("JWT_PRIVATE_KEY"))
		privKey, _ := pem.Decode(keyPem)

		var privPemBytes []byte
		if privKey.Type != "RSA PRIVATE KEY" {
			server.logger.Info("RSA private key is of the wrong type")
			rest.ERROR(w, err)
			return
		}

		privPemBytes = privKey.Bytes

		var parsedKey interface{}
		if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
			if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
				server.logger.Info("Unable to parse RSA private key")
				rest.ERROR(w, err)
				return
			}
		}

		privateKey, ok := parsedKey.(*rsa.PrivateKey)
		if !ok {
			server.logger.Info("Unable to parse RSA private key")
			rest.ERROR(w, err)
			return
		}

		err = privateKey.Validate()
		if err != nil {
			server.logger.Info("Unable to validate private key", zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		signed, err := jwt.Sign(token, jwa.RS256, privateKey)
		if err != nil {
			server.logger.Info("Unable to sign JWT with provided RSA private key", zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		body := models.JwtDO{
			Token: signed,
		}
		rest.JSON(w, http.StatusOK, body)
	}
}
