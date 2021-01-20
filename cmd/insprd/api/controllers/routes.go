package controller

import (
	"fmt"
	"net/http"
)

func (s *Server) initRoutes() {
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})
	s.Mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "App!")
	})

	s.Mux.HandleFunc("/cluster-state", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "funciton to return the cluster state!")
	})
}
