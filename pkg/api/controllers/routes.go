package controller

import (
	"net/http"

	handler "github.com/inspr/inspr/pkg/api/handlers"
)

func (s *Server) initRoutes() {
	logger.Debug("initializing Insprd server routes")
	h := handler.NewHandler(
		s.MemoryManager, s.op,
	)

	ahandler := h.NewAppHandler()
	s.Mux.HandleFunc("/apps", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			ahandler.HandleGet().JSON().Recover()(w, r)

		case http.MethodPost:
			ahandler.HandleCreate().JSON().Recover()(w, r)

		case http.MethodPut:
			ahandler.HandleUpdate().JSON().Recover()(w, r)

		case http.MethodDelete:
			ahandler.HandleDelete().JSON().Recover()(w, r)

		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	chandler := h.NewChannelHandler()
	s.Mux.HandleFunc("/channels", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			chandler.HandleGet().JSON().Recover()(w, r)

		case http.MethodPost:
			chandler.HandleCreate().JSON().Recover()(w, r)

		case http.MethodPut:
			chandler.HandleUpdate().JSON().Recover()(w, r)

		case http.MethodDelete:
			chandler.HandleDelete().JSON().Recover()(w, r)

		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	cthandler := h.NewChannelTypeHandler()
	s.Mux.HandleFunc("/channeltypes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			cthandler.HandleGet().JSON().Recover()(w, r)

		case http.MethodPost:
			cthandler.HandleCreate().JSON().Recover()(w, r)

		case http.MethodPut:
			cthandler.HandleUpdate().JSON().Recover()(w, r)

		case http.MethodDelete:
			cthandler.HandleDelete().JSON().Recover()(w, r)

		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	aliasHandler := h.NewAliasHandler()
	s.Mux.HandleFunc("/alias", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			aliasHandler.HandleGet().JSON().Recover()(w, r)

		case http.MethodPost:
			aliasHandler.HandleCreate().JSON().Recover()(w, r)

		case http.MethodPut:
			aliasHandler.HandleUpdate().JSON().Recover()(w, r)

		case http.MethodDelete:
			aliasHandler.HandleDelete().JSON().Recover()(w, r)

		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}
