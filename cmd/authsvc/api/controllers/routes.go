package controllers

import (
	"net/http"

	"inspr.dev/inspr/pkg/rest"
)

func (s *Server) initRoutes() {
	s.Mux.HandleFunc("/token", s.Tokenize().Methods(http.MethodPost))
	s.Mux.HandleFunc("/refresh", s.Refresh().Methods(http.MethodGet))
	s.Mux.HandleFunc("/init", s.HandleInit())
	s.Mux.HandleFunc("/healthz", rest.Healthz())

	// standard paths for /net/http/pprof
	rest.AttachProfiler(s.Mux)
}
