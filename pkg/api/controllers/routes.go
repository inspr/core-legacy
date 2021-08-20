package controller

import (
	"inspr.dev/inspr/pkg/rest"

	handler "inspr.dev/inspr/pkg/api/handlers"
	metabrokers "inspr.dev/inspr/pkg/meta/brokers"
)

func (s *Server) initRoutes() {

	logger.Debug("initializing Insprd server routes")
	h := handler.NewHandler(
		s.memory, s.op, s.auth,
	)

	ahandler := h.NewAppHandler()
	s.mux.Handle("/apps", rest.HandleCRUD(ahandler))

	chandler := h.NewChannelHandler()
	s.mux.Handle("/channels", rest.HandleCRUD(chandler))

	thandler := h.NewTypeHandler()
	s.mux.Handle("/types", rest.HandleCRUD(thandler))

	aliasHandler := h.NewAliasHandler()
	s.mux.Handle("/alias", rest.HandleCRUD(aliasHandler))

	brokersHandler := h.NewBrokerHandler()
	s.mux.Handle("/brokers", brokersHandler.HandleGet().Get().JSON().Validate(s.auth))
	s.mux.Handle(
		"/brokers/"+metabrokers.Kafka,
		brokersHandler.KafkaCreateHandler().Post().JSON(),
	)

	s.mux.Handle("/auth", h.TokenHandler().Validate(s.auth))
	s.mux.Handle("/refreshController", h.ControllerRefreshHandler())
	s.mux.Handle("/init", h.InitHandler())
	s.mux.Handle("/healthz", rest.Healthz())

	s.mux.Handle("/log/level", alevel)

	// standard paths for /net/http/pprof
	rest.AttachProfiler(s.mux)
}
