// Package jwtauth is responsible for implementing the auth
// methods specified in the auth folder of the inspr pkg.
package jwtauth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/rest/request"
)

// JWTauth implements the Auth interface for jwt authetication provider
type JWTauth struct {
	publicKey *rsa.PublicKey
	authURL   string
}

var logger *zap.Logger

func init() {
	logger, _ = logs.Logger(zap.Fields(zap.String("section", "authentication")))
}

// NewJWTauth takes an *rsa.PublicKey and returns an
// structure that implements the auth interface
func NewJWTauth(rsaPublicKey *rsa.PublicKey) *JWTauth {
	url, ok := os.LookupEnv("AUTH_PATH")
	if !ok {
		logger.Panic("missing AUTH_PATH environment variable")
	}
	return &JWTauth{
		publicKey: rsaPublicKey,
		authURL:   url,
	}
}

// Validate is a wrapper that checks the token of the http request and if it's
// valid, proceeds to execute the request and if it isn't valid returns an error
func (JA *JWTauth) Validate(token []byte) (*auth.Payload, []byte, error) {

	logger.Debug("parsing jwt token")
	_, err := jwt.Parse(
		token,
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, JA.publicKey),
	)
	if err != nil {
		logger.Debug("error in parsing jwt token")
		if err.Error() == errors.New(`exp not satisfied`).Error() {
			// token expired
			logger.Info("refreshing jwt token")
			newToken, err := JA.Refresh(token)
			if err != nil {
				logger.Error("error refreshing jwt token", zap.Error(err))
				return nil,
					token,
					ierrors.
						NewError().
						InternalServer().
						Message("error refreshing token: %v", err).
						Build()
			}
			token = newToken
		} else {
			return nil, token, err
		}
	}

	logger.Debug("desserializing token and acquiring credentials")
	// gets payload from token
	payload, err := auth.Desserialize(token)
	if err != nil {
		logger.Error("error desserializing token", zap.Error(err))
		return nil,
			token,
			ierrors.
				NewError().
				InternalServer().
				Message("error desserializing the payload").
				Build()
	}

	return payload, token, nil
}

// Init receives a payload and returns it in signed jwt format. Uses JWT authentication provider
func (JA *JWTauth) Init(key string, load auth.Payload) ([]byte, error) {
	logger.Debug("received initialization request")
	initDO := auth.InitDO{
		Key:     key,
		Payload: load,
	}

	client := request.NewJSONClient(JA.authURL)

	data := auth.JwtDO{}
	logger.Debug("sending initialization request to auth service", zap.String("auth-service", JA.authURL))
	err := client.Send(
		context.Background(),
		"/init",
		http.MethodPost,
		initDO,
		&data)

	if err != nil {
		logger.Error("error initializing cluster", zap.Error(err))
		err = ierrors.NewError().InternalServer().Message("error initializing cluster").InnerError(err).Build()
		return nil, err
	}

	return data.Token, nil
}

// Tokenize receives a payload and returns it in signed jwt format. Uses JWT authentication provider
func (JA *JWTauth) Tokenize(load auth.Payload) ([]byte, error) {
	logger.Debug("received tokenization request")
	client := request.NewJSONClient(JA.authURL)

	data := auth.JwtDO{}
	logger.Debug("sending request to authorization server", zap.String("auth-service", JA.authURL))
	err := client.Send(
		context.Background(),
		"/token",
		http.MethodPost,
		load,
		&data)

	if err != nil {
		logger.Error("unable to tokenize data", zap.Any("data", load), zap.String("auth-service", JA.authURL), zap.Error(err))
		err = ierrors.NewError().InternalServer().Message(err.Error()).Build()
		return nil, err
	}

	return data.Token, nil
}

// Refresh refreshes a jwt token. Uses JWT authentication provider
func (JA *JWTauth) Refresh(token []byte) ([]byte, error) {
	logger.Debug("received refresh request")
	client := request.NewClient().
		BaseURL(JA.authURL).
		Encoder(json.Marshal).
		Decoder(request.JSONDecoderGenerator).
		Header("Authorization", fmt.Sprintf("Bearer %v", string(token)))

	data := auth.JwtDO{}

	logger.Debug("sending request to authorization server", zap.String("auth-service", JA.authURL))
	err := client.Send(
		context.Background(),
		"/refresh",
		http.MethodGet,
		nil,
		&data)

	if err != nil {
		logger.Error("unable to refresh token", zap.String("auth-service", JA.authURL), zap.Error(err))
		err = ierrors.NewError().InternalServer().Message(err.Error()).Build()
		return nil, err
	}

	return data.Token, nil
}
