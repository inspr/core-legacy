package handler

import (
	"fmt"
	"net/http"
)

// Handler is an abreviation of the router function
type Handler func(w http.ResponseWriter, r *http.Request)

// Get runs function if the request method is GET
func (h Handler) Get() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			fmt.Fprint(w, "Please use a GET method for this route")
			return
		}
		h(w, r)
	}
}

// Post runs function if the request method is POST
func (h Handler) Post() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			fmt.Fprint(w, "Please use a POST method for this route")
			return
		}
		h(w, r)
	}
}

// Delete runs function if the request method is DELETE
func (h Handler) Delete() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			fmt.Fprint(w, "Please use a DELETE method for this route")
			return
		}
		h(w, r)
	}
}

// Put runs function if the request method is PUT
func (h Handler) Put() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			fmt.Fprint(w, "Please use a PUT method for this route")
			return
		}
		h(w, r)
	}
}
