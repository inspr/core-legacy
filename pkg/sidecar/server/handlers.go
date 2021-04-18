package sidecarserv

import (
	"encoding/json"
	"net/http"
	"sync"

	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/rest"
	"inspr.dev/inspr/pkg/sidecar/models"
	"go.uber.org/zap"
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

var logger *zap.Logger

// newCustomHandlers returns a struct composed of the
// Reader and Writer given in the parameters
func newCustomHandlers(l sync.Locker, r models.Reader, w models.Writer) *customHandlers {
	logger, _ = zap.NewDevelopment(
		zap.Fields(zap.String("id", environment.GetInsprAppID()), zap.String("section", "sidecar handlers")),
	)
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
	logger.Info("handling message write")
	body := models.BrokerData{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		insprError := ierrors.NewError().BadRequest().Message("couldn't parse body")
		rest.ERROR(w, insprError.Build())
		return
	}

	if !environment.IsInChannelBoundary(body.Channel, ch.OutputChannels) {
		insprError := ierrors.
			NewError().
			BadRequest().
			Message(
				"channel '%s' not found",
				body.Channel,
			)

		rest.ERROR(w, insprError.Build())
		return
	}
	logger.Info("writing message to broker", zap.String("channel", body.Channel))
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
	logger.Info("handling message read")

	body := models.BrokerData{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		insprError := ierrors.NewError().BadRequest().Message("couldn't parse body")
		rest.ERROR(w, insprError.Build())
		return
	}

	if !environment.IsInChannelBoundary(body.Channel, ch.InputChannels) {
		insprError := ierrors.
			NewError().
			BadRequest().
			Message(
				"channel '%s' not found",
				body.Channel,
			)

		rest.ERROR(w, insprError.Build())
		return
	}
	logger.Info("reading message from broker", zap.String("channel", body.Channel))

	brokerResp, err := ch.r.ReadMessage(body.Channel)
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

	if !environment.IsInChannelBoundary(body.Channel, ch.InputChannels) {
		insprError := ierrors.
			NewError().
			BadRequest().
			Message(
				"channel '%s' not found",
				body.Channel,
			)

		rest.ERROR(w, insprError.Build())
		return
	}

	if err := ch.r.Commit(body.Channel); err != nil {
		insprError := ierrors.NewError().InternalServer().InnerError(err).Message("broker's commitMessage failed")
		rest.ERROR(w, insprError.Build())
	}
}
