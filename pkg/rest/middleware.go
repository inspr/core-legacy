package rest

import "net/http"

// SetMiddlewareJSON specifies that the response content is a json
func SetMiddlewareJSON(h Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	}
}
