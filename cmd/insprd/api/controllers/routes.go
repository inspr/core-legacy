package controller

import (
	"fmt"
	"net/http"
)

func (s *Server) initRoutes() {
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
	s.Mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "App!")
	})
	// todo routes
	// create/delete/update/get routes for channels.
	// GET 		/channels 			->	returns all channels
	// GET 		/channels/new		->	returns how to make a new channel
	// POST		/channels 			->	creates a new channel
	// GET		/channels/{ref}	->	returns info about one channel in ref
	// PUT		/channels/{ref}	-> 	modifies the existing channel
	// DELETE	/channels/{ref}	->	deletes a specific channel

	// create/delete/update/get routes for apps.
	// same as channels

	// create/delete/update/get routes for channel types.
	// same as channels

}
