package controller

import (
	"net/http"

	handler "gitlab.inspr.dev/inspr/core/cmd/insprd/api/handlers"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

func (s *Server) initRoutes() {
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})

	chandler := handler.NewChannelHandler(s.MemoryManager)
	s.Mux.HandleFunc("/channels", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rest.SetMiddlewareJSON(chandler.HandleGetAllChannels())(w, r)
		case http.MethodPost:
			rest.SetMiddlewareJSON(chandler.HandleCreateChannel())(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	s.Mux.HandleFunc("/channels/info", chandler.HandleCreateInfo().Get())
	s.Mux.HandleFunc("/channels/ref", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			chandler.HandleGetChannelByRef()(w, r)
		case http.MethodPut:
			chandler.HandleUpdateChannel()(w, r)
		case http.MethodDelete:
			chandler.HandleDeleteChannel()(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	ahandler := handler.NewAppHandler(s.MemoryManager)
	s.Mux.HandleFunc("/apps", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ahandler.HandleGetAllApps()(w, r)
		case http.MethodPost:
			ahandler.HandleCreateApp()(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	s.Mux.HandleFunc("/apps/info", ahandler.HandleCreateInfo().Get())
	s.Mux.HandleFunc("/apps/ref", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ahandler.HandleGetAppByRef()(w, r)
		case http.MethodPut:
			ahandler.HandleUpdateApp()(w, r)
		case http.MethodDelete:
			ahandler.HandleDeleteApp()(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	cthandler := handler.NewChannelTypeHandler(s.MemoryManager)
	s.Mux.HandleFunc("/channeltypes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cthandler.HandleGetAllChannelTypes()(w, r)
		case http.MethodPost:
			cthandler.HandleCreateChannelType()(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	s.Mux.HandleFunc("/channeltypes/info", cthandler.HandleCreateInfo().Get())
	s.Mux.HandleFunc("/channeltypes/ref", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cthandler.HandleGetChannelTypeByRef()(w, r)
		case http.MethodPut:
			cthandler.HandleUpdateChannelType()(w, r)
		case http.MethodDelete:
			cthandler.HandleDeleteChannelType()(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}
