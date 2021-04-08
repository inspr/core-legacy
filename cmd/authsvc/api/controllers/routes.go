package controllers

import (
	"net/http"
)

func (s *Server) initRoutes() {
	s.Mux.HandleFunc("/token", s.Tokenize().Methods(http.MethodPost))
}
