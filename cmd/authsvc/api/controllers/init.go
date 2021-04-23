package controllers

import (
	"net/http"

	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
)

var initialized bool

func (server *Server) HandleInit() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		if initialized {
			rest.ERROR(w, ierrors.NewError().Message("already initialized").Build())
			return
		}
		initialized = true
		server.Tokenize()(w, r)
	}).Post().JSON().Recover()
}
