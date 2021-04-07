package rest

import (
	"errors"
	"net/http"

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

func (h Handler) Valide(auth authentication.Auth) Handler {
	return func(w http.ResponseWriter, r *http.Request) {

		// could it be that the user will send many tokens and that one of them will be valid for the operation?
		tokens := r.Header["authorization"]
		for _, token := range tokens {
			payload, err := auth.Validade(token)

			if errors.Is(err, ierrors.ExpiredToken) {

			}
		}
		// token and context are valid

		// token has expired

		// token invalid or non existant

		// valid token but not enough permissions

		// default

	}
}
