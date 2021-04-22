package controller

import (
	"github.com/inspr/inspr/pkg/rest"

	handler "github.com/inspr/inspr/pkg/api/handlers"
)

func (s *Server) initRoutes() {
	logger.Debug("initializing Insprd server routes")
	h := handler.NewHandler(
		s.MemoryManager, s.op, s.auth,
	)

	ahandler := h.NewAppHandler()
	s.Mux.Handle("/apps", rest.HandleCRUD(ahandler))

	chandler := h.NewChannelHandler()
	s.Mux.Handle("/channels", rest.HandleCRUD(chandler))

	cthandler := h.NewChannelTypeHandler()
	s.Mux.Handle("/channeltypes", rest.HandleCRUD(cthandler))

	aliasHandler := h.NewAliasHandler()
	s.Mux.Handle("/alias", rest.HandleCRUD(aliasHandler))

	s.Mux.Handle("/auth", h.TokenHandler().Validate(s.auth))
}
