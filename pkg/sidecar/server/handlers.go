package sidecarserv

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

// customHandlers is a struct that contains the handlers
//  to be used to server
type customHandlers struct {
	*Server
	insprVars *environment.InsprEnvironment
}

// newCustomHandlers returns a struct composed of the
// Reader and Writer given in the parameters
func newCustomHandlers(server *Server) *customHandlers {
	return &customHandlers{
		Server:    server,
		insprVars: environment.GetEnvironment(),
	}
}

// handles the /message route in the server
func (ch *customHandlers) writeMessageHandler(w http.ResponseWriter, r *http.Request) {
	ch.Lock()
	defer ch.Unlock()

	body := models.RequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		insprError := ierrors.NewError().BadRequest().Message("couldn't parse body")
		// TODO See if it's apropriatte to change rest.Error parameters to insprErrors
		// that implies in changing the code for the master's api
		rest.ERROR(w, http.StatusBadRequest, insprError.Build())
		return
	}

	existingChannels := strings.Split(ch.insprVars.OutputChannels, ";")

	if !existsInSlice(body.Channel, existingChannels) {
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

	body := models.RequestBody{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rest.ERROR(w, http.StatusBadRequest, err)
		return
	}

	existingChannels := strings.Split(ch.insprVars.InputChannels, ";")

	if !existsInSlice(body.Channel, existingChannels) {
		rest.ERROR(w, http.StatusBadRequest, errors.New("channel doesn't exist"))
		return
	}

	brokerResp := ch.Reader.ReadMessage(body.Channel)
	if brokerResp.Error != nil {
		rest.ERROR(w, http.StatusInternalServerError, brokerResp.Error)
		return
	}

	rest.JSON(w, http.StatusOK, brokerResp)
}

// handles the /commit route in the server
func (ch *customHandlers) commitMessageHandler(w http.ResponseWriter, r *http.Request) {
	ch.Lock()
	defer ch.Unlock()

	body := models.RequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rest.ERROR(w, http.StatusBadRequest, err)
		return
	}

	existingChannels := strings.Split(ch.insprVars.OutputChannels, ";")

	if !existsInSlice(body.Channel, existingChannels) {
		rest.ERROR(w, http.StatusBadRequest, errors.New("channel doesn't exist"))
		return
	}

	if err := ch.Reader.CommitMessage(body.Channel); err != nil {
		rest.ERROR(w, http.StatusInternalServerError, err)
	}

}

// existsInSlice checks if a channel belongs to a slice of channel
func existsInSlice(channel string, channelList []string) bool {
	exists := false
	for _, c := range channelList {
		if channel == c {
			exists = true
			break
		}
	}
	return exists
}
