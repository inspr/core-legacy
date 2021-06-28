package controller

import (
	"inspr.dev/inspr/pkg/rest"

	handler "inspr.dev/inspr/pkg/api/handlers"
	metabrokers "inspr.dev/inspr/pkg/meta/brokers"
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

	thandler := h.NewTypeHandler()
	s.Mux.Handle("/types", rest.HandleCRUD(thandler))

	aliasHandler := h.NewAliasHandler()
	s.Mux.Handle("/alias", rest.HandleCRUD(aliasHandler))

	brokersHandler := h.NewBrokerHandler()
	s.Mux.Handle("/brokers", brokersHandler.HandleGet().Get().JSON())
	s.Mux.Handle(
		"/brokers/"+metabrokers.Kafka,
		brokersHandler.KafkaCreateHandler().Post().JSON(),
	)

	s.Mux.Handle("/auth", h.TokenHandler().Validate(s.auth))
	s.Mux.Handle("/refreshController", h.ControllerRefreshHandler())
	s.Mux.Handle("/init", h.InitHandler())
}
