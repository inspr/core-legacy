package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	authentication "github.com/inspr/inspr/cmd/insprd/auth"
	"github.com/inspr/inspr/pkg/ierrors"
)

// JSON specifies in the header that the response content is a json
func (h Handler) JSON() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	}
}

// Validate handles the token validation of the http requests made, it receives an implementation of the auth interface as a parameter.
func (h Handler) Validate(auth authentication.Auth) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authorization: Bearer <token>
		headerContent := r.Header["Authorization"]

		if len(headerContent) != 1 ||
			!strings.HasPrefix(headerContent[0], "Bearer ") {
			http.Error(
				w,
				"Bad Request, expected: Authorization: <token>",
				http.StatusBadRequest,
			)
			return
		}

		token := strings.TrimPrefix(headerContent[0], "Bearer ")
		payload, newToken, err := auth.Validate([]byte(token))

		// returns the same token or a refreshed one in the header of the response
		w.Header().Add("Authorization", "Bearer "+string(newToken))

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
		requestData := struct {
			Scope string `json:"scope"`
		}{}

		// Read the content
		if r.Body != nil {
			// reads body
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			// unmarshal into scope Data
			json.Unmarshal(bodyBytes, &requestData)
			// Restore the r.Body to its original state
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		valid := false
		for _, scope := range payload.Scope {
			if strings.HasPrefix(requestData.Scope, scope) {
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
