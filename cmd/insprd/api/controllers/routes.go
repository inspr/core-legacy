package controller

import (
	"net/http"

	handler "gitlab.inspr.dev/inspr/core/cmd/insprd/api/handlers"
)

func (s *Server) initRoutes() {

	ahandler := handler.NewAppHandler(s.MemoryManager)
	s.Mux.HandleFunc("/apps", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			ahandler.HandleGetAppByRef().JSON()(w, r)

		case http.MethodPost:
			ahandler.HandleCreateApp().JSON()(w, r)

		case http.MethodPut:
			ahandler.HandleUpdateApp().JSON()(w, r)

		case http.MethodDelete:
			ahandler.HandleDeleteApp().JSON()(w, r)

		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	chandler := handler.NewChannelHandler(s.MemoryManager)
	s.Mux.HandleFunc("/channels", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			chandler.HandleGetChannelByRef().JSON()(w, r)

		case http.MethodPost:
			chandler.HandleCreateChannel().JSON()(w, r)

		case http.MethodPut:
			chandler.HandleUpdateChannel().JSON()(w, r)

		case http.MethodDelete:
			chandler.HandleDeleteChannel().JSON()(w, r)

		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	cthandler := handler.NewChannelTypeHandler(s.MemoryManager)
	s.Mux.HandleFunc("/channeltypes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			cthandler.HandleGetChannelTypeByRef().JSON()(w, r)

		case http.MethodPost:
			cthandler.HandleCreateChannelType().JSON()(w, r)

		case http.MethodPut:
			cthandler.HandleUpdateChannelType().JSON()(w, r)

		case http.MethodDelete:
			cthandler.HandleDeleteChannelType().JSON()(w, r)

		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}
