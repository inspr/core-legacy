package dapp

import (
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/sidecar"
)

type handlers struct {
	s *Server
}

func newHandlers(server *Server) *handlers {
	return &handlers{
		s: server,
	}
}

// messageHandler listens to the /message route and when it triggers
// processes the message and sends it to the app
func (h *handlers) messageHandler(w http.ResponseWriter, r *http.Request) {
	h.s.Mutex.Lock()
	defer h.s.Mutex.Unlock()

	// Add message in body to the Messages of the server
	decoder := json.NewDecoder(r.Body)
	msg := sidecar.Message{}

	if err := decoder.Decode(&msg); err != nil {
		// returns to app that there was an error
		// todo: how to do it better
		w.WriteHeader(500)
		return
	}

	h.s.Messages <- msg
}
