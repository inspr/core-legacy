package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	authentication "gitlab.inspr.dev/inspr/core/cmd/insprd/auth"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// JSON specifies in the header that the response content is a json
func (h Handler) JSON() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	}
}

func (h Handler) Validate(auth authentication.Auth) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authorization: Bearer <token>
		headerContent := r.Header["Authorization"]
		if len(headerContent) != 2 {
			http.Error(
				w,
				"Bad Request, expected: Authorization: Bearer <token>",
				http.StatusBadRequest,
			)
			return
		}

		token := headerContent[1]
		payload, _, err := auth.Validade([]byte(token))

		// error management
		if err != nil {
			// check for invalid error or non Existant
			if ierrors.HasCode(err, ierrors.InvalidToken) {
				http.Error(
					w,
					"Invalid Token",
					http.StatusUnauthorized,
				)
				return
			}

			// token expired
			if ierrors.HasCode(err, ierrors.ExpiredToken) {
				http.Error(
					w,
					"Request is OK but the token is expired",
					http.StatusOK,
				)
				return
			}

			// default error message
			http.Error(
				w,
				"Unknown error, please check token",
				http.StatusBadRequest,
			)
			return
		}

		// request scope
		data := struct {
			Scope string `json:"scope"`
		}{}

		body, _ := r.GetBody()
		json.NewDecoder(body).Decode(&data)

		valid := false
		for _, scope := range payload.Scope {
			if strings.Contains(scope, data.Scope) {
				// scope found
				valid = true
			}
		}

		// check for unauthorized error
		if !valid {
			http.Error(
				w,
				"Unauthorized to do operations in this context",
				http.StatusForbidden,
			)
			return
		}

		// token and context are valid
		h(w, r)
	}
}

// CRUDHandler handles crud requests to a given resource
type CRUDHandler interface {
	HandleCreate() Handler
	HandleDelete() Handler
	HandleUpdate() Handler
	HandleGet() Handler
}

// HandleCRUD uses a CRUDHandler to handle HTTP requests for a CRUD resource
func HandleCRUD(handler CRUDHandler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			handler.HandleGet().JSON().Recover()(w, r)

		case http.MethodPost:
			handler.HandleCreate().JSON().Recover()(w, r)

		case http.MethodPut:
			handler.HandleUpdate().JSON().Recover()(w, r)

		case http.MethodDelete:
			handler.HandleDelete().JSON().Recover()(w, r)

		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
