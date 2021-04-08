package controller

import (
	handler "gitlab.inspr.dev/inspr/core/cmd/insprd/api/handlers"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

func (s *Server) initRoutes() {
	logger.Debug("initializing Insprd server routes")
	h := handler.NewHandler(
		s.MemoryManager, s.op,
	)

	ahandler := h.NewAppHandler()
	s.Mux.Handle("/apps", rest.HandleCRUD(ahandler))

	chandler := h.NewChannelHandler()
	s.Mux.Handle("/channels", rest.HandleCRUD(chandler))

	cthandler := h.NewChannelTypeHandler()
	s.Mux.Handle("/channeltypes", rest.HandleCRUD(cthandler))

	aliasHandler := h.NewAliasHandler()
	s.Mux.Handle("/alias", rest.HandleCRUD(aliasHandler))

	s.Mux.Handle("/auth", h.TokenHandler())
}
