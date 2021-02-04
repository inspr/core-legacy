package sidecarserv

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

// customHandlers is a struct that contains the handlers
//  to be used to server
type customHandlers struct {
	*Server
}

// newCustomHandlers returns a struct composed of the
// Reader and Writer given in the parameters
func newCustomHandlers(server *Server) *customHandlers {
	return &customHandlers{server}
}

// handles the /message route in the server
func (ch *customHandlers) writeMessageHandler(w http.ResponseWriter, r *http.Request) {
	ch.Lock()
	defer ch.Unlock()

	decoder := json.NewDecoder(r.Body)
	body := models.RequestBody{}

	if err := decoder.Decode(&body); err != nil {
		rest.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// todo use the environment method when is ready
	existingChannels := strings.Split(environment.GetEnvironment().OutputChannels, ";")

	// todo separate function
	exists := false
	for _, envChan := range existingChannels {
		if body.Channel == envChan {
			exists = true
			break
		}
	}

	if !exists {
		rest.ERROR(w, http.StatusBadRequest, errors.New("channel doesn't exist"))
		return
	}

	if err := ch.Writer.WriteMessage(body.Channel, body.Message); err != nil {
		rest.ERROR(w, http.StatusInternalServerError, err)
	}
}

// handles the /message route in the server
func (ch *customHandlers) readMessageHandler(w http.ResponseWriter, r *http.Request) {
	ch.Lock()
	defer ch.Unlock()

	decoder := json.NewDecoder(r.Body)
	body := models.RequestBody{}

	if err := decoder.Decode(&body); err != nil {
		rest.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// todo use the environment method when is ready
	existingChannels := strings.Split(environment.GetEnvironment().InputChannels, ";")

	// todo make it not hideous
	exists := false
	for _, envChan := range existingChannels {
		if body.Channel == envChan {
			exists = true
			break
		}
	}
	if !exists {
		rest.ERROR(w, http.StatusBadRequest, errors.New("channel doesn't exist"))
		return
	}

	msg, err := ch.Reader.ReadMessage(body.Channel)
	if err != nil {
		rest.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// todo return in body the msg
	rest.JSON(w, http.StatusOK, msg)
}

// handles the /commit route in the server
func (ch *customHandlers) commitMessageHandler(w http.ResponseWriter, r *http.Request) {
	ch.Lock()
	defer ch.Unlock()

	decoder := json.NewDecoder(r.Body)
	body := models.RequestBody{}

	if err := decoder.Decode(&body); err != nil {
		rest.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// todo use the environment method when is ready
	existingChannels := strings.Split(environment.GetEnvironment().OutputChannels, ";")

	// todo make it not hideous
	exists := false
	for _, envChan := range existingChannels {
		if body.Channel == envChan {
			exists = true
			break
		}
	}
	if !exists {
		rest.ERROR(w, http.StatusBadRequest, errors.New("channel doesn't exist"))
		return
	}

	if err := ch.Reader.CommitMessage(body.Channel); err != nil {
		rest.ERROR(w, http.StatusInternalServerError, err)
	}

}
