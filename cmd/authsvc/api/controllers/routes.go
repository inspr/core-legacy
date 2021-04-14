package controllers

import (
	"net/http"
)

func (s *Server) initRoutes() {
	s.Mux.HandleFunc("/token", s.Tokenize().Methods(http.MethodPost))
	s.Mux.HandleFunc("/refresh", s.Refresh().Methods(http.MethodPost))
}
