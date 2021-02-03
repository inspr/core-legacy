package sidecarserv

import (
	"bytes"
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

// handles the /message route in the server
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

// handles the /commit route in the server
func (h *handlers) commitHandler(w http.ResponseWriter, r *http.Request) {
	h.s.Mutex.Lock()
	commitErrs := make([]error, 0)
	// Commits all the messages to the writeAddr
	for len(h.s.Messages) > 0 {
		msg := <-h.s.Messages

		reqBody, err := json.Marshal(msg)
		if err != nil {
			commitErrs = append(commitErrs, err)
			continue
		}

		response, err := http.Post(h.s.sendAddr, "", bytes.NewReader(reqBody))
		if response.StatusCode != http.StatusOK {
			// todo: see how the err comes back in request
			// is it in the 'err' of the body? if so just unmarshal
			commitErrs = append(commitErrs)
			continue
		}
	}
	h.s.Mutex.Unlock()
}
