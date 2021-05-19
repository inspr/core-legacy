package sidecarserv

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/sidecar_old/models"
	"go.uber.org/zap"
)

var logger *zap.Logger

const maxBrokerRetries = 5

type response struct {
	Status string
}

var (
	writeMessageErr = ierrors.NewError().InternalServer().Message("broker's writeMessage failed")
	decodingErr     = ierrors.NewError().BadRequest().Message("couldn't parse body")
)

func init() {
	logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "loadbalencer-sidecar")))
}

// handles the /message route in the server
func (s *Server) writeMessageHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling message write")

		body := models.BrokerData{}
		channel := strings.TrimPrefix(r.URL.Path, "/")

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			rest.ERROR(w, decodingErr.InnerError(err).Build())
			return
		}

		if !environment.OutputChannnelList().Contains(channel) {
			insprError := ierrors.
				NewError().
				BadRequest().
				Message(
					"channel '%s' not found",
					channel,
				)

			rest.ERROR(w, insprError.Build())
			return
		}

		// get env var nisso que vai dar bom channel+"_RESOLVED"

		logger.Info("sending message to broker",
			zap.String("broker", channel),
			zap.String("channel", channel))

		if err := s.Writer.WriteMessage(channel, body.Message); err != nil {
			rest.ERROR(w, writeMessageErr.InnerError(err).Build())
			return
		}
		rest.JSON(w, 200, struct{ Status string }{"OK"})
	}
}

func (s *Server) writeWithRetry(ctx context.Context, channel string, data interface{}) (resp response, err error) {
	for i := 0; ; i++ {
		err = s.client.Send(ctx, "/"+channel, http.MethodPost, data, &resp)
		if err != nil {
			if i == maxBrokerRetries {
				return
			}
			continue
		}
		return
	}
}

func (s *Server) readWithRetry(ctx context.Context, channel string) (brokerResp models.BrokerData, err error) {
	for i := 0; ; i++ {
		brokerResp, err = s.Reader.ReadMessage(ctx, channel)
		if err != nil {
			if i == maxBrokerRetries {
				return
			}
			continue
		}
		return
	}
}

func (s *Server) channelReadMessageRoutine(ctx context.Context, channel string) error {

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var err error
			var brokerResp models.BrokerData

			brokerResp, err = s.readWithRetry(ctx, channel)
			if err != nil {
				return err
			}

			fmt.Println("trying to send requess")

			resp, err := s.writeWithRetry(ctx, channel, brokerResp)
			if err != nil || resp.Status != "OK" {
				return err
			}
			s.Reader.Commit(ctx, channel)
		}
	}
}

func (s *Server) readMessageRoutine(ctx context.Context) error {
	s.runningRead = true
	defer func() { s.runningRead = false }()

	errch := make(chan error)
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, channel := range environment.InputChannelList() {
		go func(routeChan string) { errch <- s.channelReadMessageRoutine(newCtx, routeChan) }(channel)
	}

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}

}

// Close closes the server connection
func (s *Server) Close() {
	s.cancel()
}
