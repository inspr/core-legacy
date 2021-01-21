/*
Package rest contains the functions
that make it easier to manager api
handler functions
*/
package rest

import (
	"net/http"
)

// Handler is an abreviation of the api router function
type Handler func(w http.ResponseWriter, r *http.Request)

// Get runs function if the request method is GET
func (h Handler) Get() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			// todo: review if below should be done
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

// Post runs function if the request method is POST
func (h Handler) Post() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			// todo: review if below should be done
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

// Delete runs function if the request method is DELETE
func (h Handler) Delete() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			// todo: review if below should be done
			w.Header().Set("Allow", http.MethodDelete)
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

// Put runs function if the request method is PUT
func (h Handler) Put() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			// todo: review if below should be done
			w.Header().Set("Allow", http.MethodPut)
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}
