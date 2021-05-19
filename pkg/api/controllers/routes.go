package controller

import (
	"github.com/inspr/inspr/pkg/rest"

	handler "github.com/inspr/inspr/pkg/api/handlers"
)

func (s *Server) initRoutes() {
	logger.Debug("initializing Insprd server routes")
	h := handler.NewHandler(
		s.MemoryManager, s.op, s.auth, s.BrokerManager,
	)

	ahandler := h.NewAppHandler()
	s.Mux.Handle("/apps", rest.HandleCRUD(ahandler))

	chandler := h.NewChannelHandler()
	s.Mux.Handle("/channels", rest.HandleCRUD(chandler))

	cthandler := h.NewTypeHandler()
	s.Mux.Handle("/types", rest.HandleCRUD(cthandler))

	aliasHandler := h.NewAliasHandler()
	s.Mux.Handle("/alias", rest.HandleCRUD(aliasHandler))

	brokersHandler := h.NewBrokerHandler()
	s.Mux.Handle("/brokers", brokersHandler.HandleGet())

	s.Mux.Handle("/auth", h.TokenHandler().Validate(s.auth))
	s.Mux.Handle("/refreshController", h.ControllerRefreshHandler())
	s.Mux.Handle("/init", h.InitHandler())
}
