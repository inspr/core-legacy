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
			// rest.SetMiddlewareJSON(chandler.HandleGetAllChannels())(w, r)
		case http.MethodPost:
			rest.SetMiddlewareJSON(chandler.HandleCreateChannel())(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	s.Mux.HandleFunc("/channels/info",
		rest.SetMiddlewareJSON(chandler.HandleCreateInfo()).Get(),
	)
	s.Mux.HandleFunc("/channels/ref", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rest.SetMiddlewareJSON(rest.SetMiddlewareJSON(chandler.HandleGetChannelByRef()))(w, r)
		case http.MethodPut:
			rest.SetMiddlewareJSON(chandler.HandleUpdateChannel())(w, r)
		case http.MethodDelete:
			rest.SetMiddlewareJSON(chandler.HandleDeleteChannel())(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	ahandler := handler.NewAppHandler(s.MemoryManager)
	s.Mux.HandleFunc("/apps", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// rest.SetMiddlewareJSON(ahandler.HandleGetAllApps())(w, r)
		case http.MethodPost:
			rest.SetMiddlewareJSON(ahandler.HandleCreateApp())(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	s.Mux.HandleFunc("/apps/info",
		rest.SetMiddlewareJSON(ahandler.HandleCreateInfo()).Get(),
	)
	s.Mux.HandleFunc("/apps/ref", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rest.SetMiddlewareJSON(ahandler.HandleGetAppByRef())(w, r)
		case http.MethodPut:
			rest.SetMiddlewareJSON(ahandler.HandleUpdateApp())(w, r)
		case http.MethodDelete:
			rest.SetMiddlewareJSON(ahandler.HandleDeleteApp())(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	cthandler := handler.NewChannelTypeHandler(s.MemoryManager)
	s.Mux.HandleFunc("/channeltypes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// rest.SetMiddlewareJSON(cthandler.HandleGetAllChannelTypes())(w, r)
		case http.MethodPost:
			rest.SetMiddlewareJSON(cthandler.HandleCreateChannelType())(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	s.Mux.HandleFunc("/channeltypes/info",
		rest.SetMiddlewareJSON(cthandler.HandleCreateInfo()).Get(),
	)
	s.Mux.HandleFunc("/channeltypes/ref", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rest.SetMiddlewareJSON(cthandler.HandleGetChannelTypeByRef())(w, r)
		case http.MethodPut:
			rest.SetMiddlewareJSON(cthandler.HandleUpdateChannelType())(w, r)
		case http.MethodDelete:
			rest.SetMiddlewareJSON(cthandler.HandleDeleteChannelType())(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}
