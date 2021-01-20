package controller

import (
	"net/http"

	handler "gitlab.inspr.dev/inspr/core/cmd/insprd/api/handlers"
)

func (s *Server) initRoutes() {
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
	// todo routes
	// create/delete/update/get routes for channels.
	// GET 		/channels 			->	returns all channels
	// GET 		/channels/new		->	returns how to make a new channel
	// POST		/channels 			->	creates a new channel
	// GET		/channels/{ref}	->	returns info about one channel in ref
	// PUT		/channels/{ref}	-> 	modifies the existing channel
	// DELETE	/channels/{ref}	->	deletes a specific channel
	s.Mux.HandleFunc("/channels", handler.ChannelHandler.GetAllChannels().Get())
	s.Mux.HandleFunc("/channels/new", handler.ChannelHandler.CreateInfo().Get())
	s.Mux.HandleFunc("/channels", handler.ChannelHandler.CreateChannel().Post())
	s.Mux.HandleFunc("/channels/{ref}", handler.ChannelHandler.GetChannelByRef().Get())
	s.Mux.HandleFunc("/channels/{ref}", handler.ChannelHandler.UpdateChannel().Put())
	s.Mux.HandleFunc("/channels/{ref}", handler.ChannelHandler.DeleteChannel().Delete())

	// create/delete/update/get routes for apps.
	s.Mux.HandleFunc("/apps", handler.ChannelHandler.GetAllChannels().Get())
	s.Mux.HandleFunc("/apps/new", handler.ChannelHandler.CreateInfo().Get())
	s.Mux.HandleFunc("/apps", handler.ChannelHandler.CreateChannel().Post())
	s.Mux.HandleFunc("/apps/{ref}", handler.ChannelHandler.GetChannelByRef().Get())
	s.Mux.HandleFunc("/apps/{ref}", handler.ChannelHandler.UpdateChannel().Put())
	s.Mux.HandleFunc("/apps/{ref}", handler.ChannelHandler.DeleteChannel().Delete())

	// create/delete/update/get routes for channel types.
	s.Mux.HandleFunc("/channeltypes", handler.ChannelHandler.GetAllChannels().Get())
	s.Mux.HandleFunc("/channeltypes/new", handler.ChannelHandler.CreateInfo().Get())
	s.Mux.HandleFunc("/channeltypes", handler.ChannelHandler.CreateChannel().Post())
	s.Mux.HandleFunc("/channeltypes/{ref}", handler.ChannelHandler.GetChannelByRef().Get())
	s.Mux.HandleFunc("/channeltypes/{ref}", handler.ChannelHandler.UpdateChannel().Put())
	s.Mux.HandleFunc("/channeltypes/{ref}", handler.ChannelHandler.DeleteChannel().Delete())

}
