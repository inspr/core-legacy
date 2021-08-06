package rest

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/utils"
)

// CRUDHandler handles crud requests to a given resource
type CRUDHandler interface {
	HandleCreate() Handler
	HandleDelete() Handler
	HandleUpdate() Handler
	HandleGet() Handler
	GetAuth() auth.Auth
	GetCancel() func()
}

// HandleCRUD uses a CRUDHandler to handle HTTP requests for a CRUD resource
func HandleCRUD(handler CRUDHandler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			handler.
				HandleGet().
				Validate(handler.GetAuth()).
				JSON().
				Recover(handler.GetCancel())(w, r)

		case http.MethodPost:
			handler.
				HandleCreate().
				Validate(handler.GetAuth()).
				JSON().
				Recover(handler.GetCancel())(w, r)

		case http.MethodPut:
			handler.
				HandleUpdate().
				Validate(handler.GetAuth()).
				JSON().
				Recover(handler.GetCancel())(w, r)

		case http.MethodDelete:
			handler.
				HandleDelete().
				Validate(handler.GetAuth()).
				JSON().
				Recover(handler.GetCancel())(w, r)

		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

// JSON specifies in the header that the response content is a json
func (h Handler) JSON() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	}
}

// Validate handles the token validation of the http requests made, it receives an implementation of the auth interface as a parameter.
func (h Handler) Validate(auth auth.Auth) Handler {
	logger, _ := logs.Logger(zap.Fields(zap.String("section", "api"), zap.String("subSection", "authorization-middleware")))
	return func(w http.ResponseWriter, r *http.Request) {
		// Authorization: Bearer <token>
		headerContent := r.Header["Authorization"]
		logger.Info("validating request")
		if (len(headerContent) == 0) ||
			(!strings.HasPrefix(headerContent[0], "Bearer ")) {
			logger.Info("invalid token received")
			ERROR(w, ierrors.New("invalid token format").Unauthorized())
			return
		}

		token := strings.TrimPrefix(headerContent[0], "Bearer ")
		payload, newToken, err := auth.Validate([]byte(token))
		logger.Debug("payload after validation")

		// returns the same token or a refreshed one in the header of the response
		w.Header().Add("Authorization", "Bearer "+string(newToken))

		// error management
		if err != nil {
			// check for invalid error or non existent
			if ierrors.HasCode(err, ierrors.InvalidToken) {
				logger.Info("invalid token received, refusing request")
				ERROR(w, ierrors.New("invalid token").Unauthorized())
				return
			}

			// default error message
			ERROR(w, ierrors.New(err))
			return
		}

		// used for checking scope authorization
		reqScopes := r.Header[HeaderScopeKey]
		logger.Debug("payload permissions", zap.Any("permissions", payload.Permissions))

		// used for checking permissions
		operation := getOperation(r)
		target := getTarget(r)
		perm := operation + ":" + target
		logger.Info("validating permissions for request", zap.String("operation", operation), zap.String("target", target))

		for scope := range payload.Permissions {
			logger.Debug("checking permissions for scope", zap.String("token-scope", scope))

			// usually the request will one have one scope
			for _, rs := range reqScopes {
				logger.Debug("comparing scope with token scope", zap.String("request-scope", rs))

				if strings.HasPrefix(rs, scope) &&
					utils.Includes(payload.Permissions[scope], perm) {
					logger.Info("permission granted for request", zap.String("request-scope", rs), zap.String("token-scope", scope), zap.String("request", perm))
					h(w, r)
					return
				}
			}
		}

		// there were no valid operations
		logger.Info("insufficient credentials, refusing request", zap.String("requested-permission", perm), zap.Strings("requested-scopes", reqScopes))
		ERROR(
			w,
			ierrors.New(
				"not enought permissions to perform request",
			).Forbidden(),
		)
	}
}

// getOperation returns the operation being done by the Request in the cluster
// get, create, update, delete.
func getOperation(r *http.Request) string {
	// some methods represent their own operation
	operation, ok := operationTranslator[r.Method]
	if !ok {
		return r.Method
	}
	return operation
}

// getTarget isolates the area that is being requested, for example the request
// URL is https://example.org:8000/channels, the getTarget removes the base of the url
// and some unnecessary '/' and returns only 'channels'
func getTarget(r *http.Request) string {
	route := strings.TrimSuffix(r.URL.Path, "/")
	route = strings.TrimPrefix(route, r.URL.Host)
	route = strings.TrimPrefix(route, "/")

	// some constant values differ from the url name used
	target, ok := routeTranslator[route]
	if !ok {
		return route
	}
	return target
}
