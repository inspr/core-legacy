package rest

import (
	"log"
	"net/http"
	"strings"

	"github.com/inspr/inspr/pkg/auth"
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
func (h Handler) Validate(auth auth.Auth) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authorization: Bearer <token>
		headerContent := r.Header["Authorization"]
		log.Println("validating")
		log.Printf("headerContent = %+v\n", headerContent)
		if (len(headerContent) == 0) ||
			(!strings.HasPrefix(headerContent[0], "Bearer ")) {

			ERROR(w, ierrors.NewError().Unauthorized().Message("invalid token format").Build())
			return
		}

		token := strings.TrimPrefix(headerContent[0], "Bearer ")
		payload, newToken, err := auth.Validate([]byte(token))
		log.Printf("payload = %+v\n", payload)
		log.Printf("string(newToken) = %+v\n", string(newToken))

		// returns the same token or a refreshed one in the header of the response
		w.Header().Add("Authorization", "Bearer "+string(newToken))

		// error management
		if err != nil {
			// check for invalid error or non Existant
			if ierrors.HasCode(err, ierrors.InvalidToken) {

				ERROR(w, ierrors.NewError().Unauthorized().Message("invalid token").Build())
				return
			}

			// default error message
			ERROR(w, ierrors.NewError().Message(err.Error()).Build())
			return
		}

		reqScopes := r.Header[HeaderScopeKey]
		valid := false
		log.Printf("payload.Permissions = %+v\n", payload.Permissions)

		for scope := range payload.Permissions {
			log.Printf("permission-scope = %+v\n", scope)
			for _, rs := range reqScopes {
				log.Printf("request-scope = %+v\n", rs)
				if strings.HasPrefix(rs, scope) {
					// scope found
					valid = true
				}
			}
		}

		// check for unauthorized error
		if !valid {
			ERROR(w, ierrors.NewError().Forbidden().Message("not enought permissions to perform request").Build())
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
	GetAuth() auth.Auth
}

// HandleCRUD uses a CRUDHandler to handle HTTP requests for a CRUD resource
func HandleCRUD(handler CRUDHandler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			handler.HandleGet().Validate(handler.GetAuth()).JSON().Recover()(w, r)

		case http.MethodPost:
			handler.HandleCreate().Validate(handler.GetAuth()).JSON().Recover()(w, r)

		case http.MethodPut:
			handler.HandleUpdate().Validate(handler.GetAuth()).JSON().Recover()(w, r)

		case http.MethodDelete:
			handler.HandleDelete().Validate(handler.GetAuth()).JSON().Recover()(w, r)

		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
