package sidecarserv

import (
	"encoding/json"
	"net/http"
	"sync"

	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

// customHandlers is a struct that contains the handlers
//  to be used to server
type customHandlers struct {
	sync.Locker
	r              models.Reader
	w              models.Writer
	InputChannels  string
	OutputChannels string
}

// newCustomHandlers returns a struct composed of the
// Reader and Writer given in the parameters
func newCustomHandlers(l sync.Locker, r models.Reader, w models.Writer) *customHandlers {
	return &customHandlers{
		Locker:         l,
		r:              r,
		w:              w,
		InputChannels:  environment.GetInputChannels(),
		OutputChannels: environment.GetOutputChannels(),
	}
}

// handles the /message route in the server
func (ch *customHandlers) writeMessageHandler(w http.ResponseWriter, r *http.Request) {
	ch.Lock()
	defer ch.Unlock()

	body := models.BrokerData{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		insprError := ierrors.NewError().BadRequest().Message("couldn't parse body")
		rest.ERROR(w, insprError.Build())
		return
	}

	if !environment.IsInOutputChannel(body.Channel, ch.OutputChannels) {
		insprError := ierrors.NewError().BadRequest().Message("channel not found")
		rest.ERROR(w, insprError.Build())
		return
	}

	if err := ch.w.WriteMessage(body.Channel, body.Message.Data); err != nil {
		insprError := ierrors.NewError().InternalServer().InnerError(err).Message("broker's writeMessage failed")
		rest.ERROR(w, insprError.Build())
		return
	}
}

// handles the /message route in the server
func (ch *customHandlers) readMessageHandler(w http.ResponseWriter, r *http.Request) {
	ch.Lock()
	defer ch.Unlock()

	body := models.BrokerData{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		insprError := ierrors.NewError().BadRequest().Message("couldn't parse body")
		rest.ERROR(w, insprError.Build())
		return
	}

	if !environment.IsInInputChannel(body.Channel, ch.InputChannels) {
		insprError := ierrors.NewError().BadRequest().Message("channel not found")
		rest.ERROR(w, insprError.Build())
		return
	}

	brokerResp, err := ch.r.ReadMessage()
	if err != nil {
		insprError := ierrors.NewError().InternalServer().InnerError(err).Message("broker's ReadMessage returned an error")
		rest.ERROR(w, insprError.Build())
		return
	}

	rest.JSON(w, http.StatusOK, brokerResp)
}

// handles the /commit route in the server
func (ch *customHandlers) commitMessageHandler(w http.ResponseWriter, r *http.Request) {
	ch.Lock()
	defer ch.Unlock()

	body := models.BrokerData{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		insprError := ierrors.NewError().BadRequest().Message("couldn't parse body")
		rest.ERROR(w, insprError.Build())
		return
	}

	if !environment.IsInInputChannel(body.Channel, ch.InputChannels) {
		insprError := ierrors.NewError().BadRequest().Message("channel not found")
		rest.ERROR(w, insprError.Build())
		return
	}

	if err := ch.r.CommitMessage(); err != nil {
		insprError := ierrors.NewError().InternalServer().InnerError(err).Message("broker's commitMessage failed")
		rest.ERROR(w, insprError.Build())
	}

}
