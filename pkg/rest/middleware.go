package rest

import (
	b64 "encoding/base64"
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
		// TODO REWRITE THIS WITH THE []BYTE RETURN
		payload, err := auth.Validade(token)

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

		// payload data
		var pl authentication.Payload
		payloadBytes, _ := b64.StdEncoding.DecodeString(payload)
		json.Unmarshal(payloadBytes, &pl)

		// TODO rewrite this to get from the HEADER
		// request scope
		data := struct {
			Scope string `json:"scope"`
		}{}
		// TODO this will be reading the body, work on a copy
		json.NewDecoder(r.Body).Decode(&data)

		for _, scope := range pl.Scope {
			if strings.Contains(scope, data.Scope) {
				// scope found
			}
		}
		// check for unauthorized error
		// if ierrors.HasCode(err, ierrors.Unauthorized) {
		// 	http.Error(
		// 		w,
		// 		"Unauthorized to do operations in this context",
		// 		http.StatusForbidden,
		// 	)
		// 	return
		// }

		// token and context are valid
		h(w, r)
	}
}
