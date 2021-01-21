package controller

import (
	"net/http"

	handler "gitlab.inspr.dev/inspr/core/cmd/insprd/api/handlers"
)

// todo routes
// GET 		/channels 			->	returns all channels
// GET 		/channels/new		->	returns how to make a new channel
// POST		/channels 			->	creates a new channel
// GET		/channels/ref	->	returns info about one channel in ref
// PUT		/channels/ref	-> 	modifies the existing channel
// DELETE	/channels/ref	->	deletes a specific channel

func (s *Server) initRoutes() {
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
	s.Mux.HandleFunc("/channels", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ChannelHandler.GetAllChannels()(w, r)
		case http.MethodPost:
			handler.ChannelHandler.CreateChannel()(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	s.Mux.HandleFunc("/channels/info", handler.ChannelHandler.CreateInfo().Get())
	s.Mux.HandleFunc("/channels/ref", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ChannelHandler.GetChannelByRef()(w, r)
		case http.MethodPut:
			handler.ChannelHandler.UpdateChannel()(w, r)
		case http.MethodDelete:
			handler.ChannelHandler.DeleteChannel()(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	s.Mux.HandleFunc("/apps", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.AppHandler.GetAllApps()(w, r)
		case http.MethodPost:
			handler.AppHandler.CreateApp()(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	s.Mux.HandleFunc("/apps/info", handler.AppHandler.CreateInfo().Get())
	s.Mux.HandleFunc("/apps/ref", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.AppHandler.GetAppByRef()(w, r)
		case http.MethodPut:
			handler.AppHandler.UpdateApp()(w, r)
		case http.MethodDelete:
			handler.AppHandler.DeleteApp()(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	s.Mux.HandleFunc("/channeltypes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ChannelTypeHandler.GetAllChannelTypes()(w, r)
		case http.MethodPost:
			handler.ChannelTypeHandler.CreateChannelType()(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
	s.Mux.HandleFunc("/channeltypes/info", handler.ChannelTypeHandler.CreateInfo().Get())
	s.Mux.HandleFunc("/channeltypes/ref", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ChannelTypeHandler.GetChannelTypeByRef()(w, r)
		case http.MethodPut:
			handler.ChannelTypeHandler.UpdateChannelType()(w, r)
		case http.MethodDelete:
			handler.ChannelTypeHandler.DeleteChannelType()(w, r)
		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	// // create/delete/update/get routes for channel types.
	// s.Mux.HandleFunc("/channeltypes", handler.ChannelHandler.GetAllChannels().Get())
	// s.Mux.HandleFunc("/channeltypes/new", handler.ChannelHandler.CreateInfo().Get())
	// s.Mux.HandleFunc("/channeltypes", handler.ChannelHandler.CreateChannel().Post())
	// s.Mux.HandleFunc("/channeltypes/{ref}", handler.ChannelHandler.GetChannelByRef().Get())
	// s.Mux.HandleFunc("/channeltypes/{ref}", handler.ChannelHandler.UpdateChannel().Put())
	// s.Mux.HandleFunc("/channeltypes/{ref}", handler.ChannelHandler.DeleteChannel().Delete())

}
