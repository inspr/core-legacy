package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/rest"
)

func init() {
	tokenLogger = logger.With(zap.String("subSection", "token"))
}

var tokenLogger *zap.Logger

// ControllerRefreshHandler handles requests for token refresing on inspr controllers on Insprd
func (h *Handler) ControllerRefreshHandler() rest.Handler {
	l := tokenLogger.With(zap.String("operation", "refresh"))
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		l.Info("received refresh token request")
		received := auth.ResfreshDO{}
		l.Debug("decoding refresh body")
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			l.Error("unable to decode body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		// this is the path to the app
		appQuery := string(received.RefreshToken)

		l.Debug("querying app to get its credentials", zap.String("app-query", appQuery))
		app, err := h.Memory.Tree().Apps().Get(appQuery)
		if err != nil {
			l.Error("error finding dApp from request", zap.String("app-query", appQuery))
			rest.ERROR(w, err)
			return
		}
		// refresh the payload with the current permissions of the dApp
		payload := auth.Payload{
			UID: app.Meta.UUID,
			Permissions: map[string][]string{
				app.Spec.Auth.Scope: app.Spec.Auth.Permissions,
			},
			Refresh:    []byte(appQuery),
			RefreshURL: fmt.Sprintf("%v/refreshController", os.Getenv("INSPR_INSPRD_ADDRESS")),
		}
		l.Debug("sucessfully refreshed token")
		rest.JSON(w, 200, payload)

	}).Recover().Post().JSON()
}

// TokenHandler handles requests for token creation on Insprd
func (h *Handler) TokenHandler() rest.Handler {
	l := tokenLogger.With(zap.String("operation", "tokenization"))
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		l.Info("received token creation request")
		var load auth.Payload
		l.Debug("decoding body")
		err := json.NewDecoder(r.Body).Decode(&load)
		if err != nil {
			l.Error("error decoding body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		l.Debug("sending tokenization request")
		token, err := h.Auth.Tokenize(load)
		if err != nil {
			l.Error("error tokenizing payload", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		l.Debug("successfully applied tokenization")
		rest.JSON(w, http.StatusOK, auth.JwtDO{
			Token: token,
		})

	}).Recover().Post().JSON()
}

// InitHandler handles requests for cluster auth initialization
func (h *Handler) InitHandler() rest.Handler {
	l := tokenLogger.With(zap.String("operation", "initialization"))
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		l.Info("received initialization request")
		l.Debug("decoding body")
		decoder := json.NewDecoder(r.Body)
		res := struct{ Key string }{}
		err := decoder.Decode(&res)
		if err != nil {
			l.Error("error decoding body", zap.Error(err))
			rest.ERROR(w, err)
		}
		load := auth.Payload{
			RefreshURL:  os.Getenv("REFRESH_URL"),
			Permissions: map[string][]string{"": {auth.CreateToken}},
		}
		l.Debug("sending request to auth for initialization")
		token, err := h.Auth.Init(res.Key, load)
		if err != nil {
			l.Error("error authenticating token", zap.Error(err))
			rest.ERROR(w, ierrors.NewError().InternalServer().Message("unable to authenticate token").InnerError(err).Build())
			return
		}
		l.Debug("successfully initialized cluster")
		rest.JSON(w, http.StatusOK, auth.JwtDO{
			Token: token,
		})
	}).Recover().Post().JSON()
}
