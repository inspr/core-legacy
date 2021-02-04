package sidecarserv

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/rest"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

// handles the /message route in the server
func (s *Server) writeMessageHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	decoder := json.NewDecoder(r.Body)
	body := struct {
		msg     models.Message `json:"msg"`
		channel string         `json:"channel"`
	}{}

	if err := decoder.Decode(&body); err != nil {
		rest.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// todo use the environment method when is ready
	existingChannels := []string{"environment.GetEnvironment().OutputChannels"}

	// todo separate function
	exists := false
	for _, envChan := range existingChannels {
		if body.channel == envChan {
			exists = true
			break
		}
	}

	if !exists {
		rest.ERROR(w, http.StatusBadRequest, errors.New("channel doesn't exist"))
		return
	}

	if err := s.Writer.WriteMessage(body.channel, body.msg); err != nil {
		rest.ERROR(w, http.StatusInternalServerError, err)
		return
	}
}

// handles the /message route in the server
func (s *Server) readMessageHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	decoder := json.NewDecoder(r.Body)
	body := struct {
		channel string `json:"channel"`
	}{}

	if err := decoder.Decode(&body); err != nil {
		rest.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// todo use the environment method when is ready
	existingChannels := []string{"environment.GetEnvironment().InputChannels"}

	// todo make it not hideous
	exists := false
	for _, envChan := range existingChannels {
		if body.channel == envChan {
			exists = true
			break
		}
	}
	if !exists {
		rest.ERROR(w, http.StatusBadRequest, errors.New("channel doesn't exist"))
		return
	}

	msg, err := s.Reader.ReadMessage(body.channel)
	if err != nil {
		rest.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// todo return in body the msg
	rest.JSON(w, http.StatusOK, msg)
}

// handles the /commit route in the server
func (s *Server) commitMessageHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	decoder := json.NewDecoder(r.Body)
	body := struct {
		channel string `json:"channel"`
	}{}

	if err := decoder.Decode(&body); err != nil {
		rest.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// todo use the environment method when is ready
	existingChannels := []string{"environment.GetEnvironment().OutputChannels"}

	// todo make it not hideous
	exists := false
	for _, envChan := range existingChannels {
		if body.channel == envChan {
			exists = true
			break
		}
	}
	if !exists {
		rest.ERROR(w, http.StatusBadRequest, errors.New("channel doesn't exist"))
		return
	}

	if err := s.Reader.CommitMessage(body.channel); err != nil {
		rest.ERROR(w, http.StatusInternalServerError, err)
	}

}
