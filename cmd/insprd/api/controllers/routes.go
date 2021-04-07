package controller

import (
	handler "gitlab.inspr.dev/inspr/core/cmd/insprd/api/handlers"
)

func (s *Server) initRoutes() {
	h := handler.NewHandler(
		s.MemoryManager, s.op,
	)

	ahandler := h.NewAppHandler()
	s.Mux.Handle("/apps", ahandler.Serve())

	chandler := h.NewChannelHandler()
	s.Mux.Handle("/channels", chandler.Serve())

	cthandler := h.NewChannelTypeHandler()
	s.Mux.Handle("/channeltypes", cthandler.Serve())

	aliasHandler := h.NewAliasHandler()
	s.Mux.HandleFunc("/alias", aliasHandler.Serve())
}
