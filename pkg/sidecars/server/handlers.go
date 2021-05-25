package sidecarserv

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

const maxBrokerRetries = 5

var (
	writeMessageErr = ierrors.NewError().InternalServer().Message("broker's writeMessage failed")
)

// handles the /message route in the server
func (s *Server) writeMessageHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling message write")

		channel := strings.TrimPrefix(r.URL.Path, "/")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		if !environment.OutputBrokerChannnels(s.broker).Contains(channel) { // OutputChannnelList must be checked for obtaining the right list
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

		logger.Info("writing message to broker", zap.String("channel", channel))
		if err := s.Writer.WriteMessage(channel, body); err != nil {
			rest.ERROR(w, writeMessageErr.InnerError(err).Build())
			return
		}
		rest.JSON(w, 200, struct{ Status string }{"OK"})
	}
}

func (s *Server) writeWithRetry(ctx context.Context, channel string, data []byte) (resp response, err error) {
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

func (s *Server) readWithRetry(ctx context.Context, channel string) (brokerMsg []byte, err error) {
	for i := 0; ; i++ {
		brokerMsg, err = s.Reader.ReadMessage(ctx, channel)
		if err != nil {
			if i == maxBrokerRetries {
				return
			}
			continue
		}
		return
	}
}

type response struct {
	Status string
}

func (s *Server) channelReadMessageRoutine(ctx context.Context, channel string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var err error
			var brokerMsg []byte

			brokerMsg, err = s.readWithRetry(ctx, channel)
			if err != nil {
				return err
			}

			logger.Debug("trying to send request to loadbalancer")

			resp, err := s.writeWithRetry(ctx, channel, brokerMsg)
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

	for _, channel := range environment.InputBrokerChannels(s.broker) { // InputChannelList retorna todods os canais de input do node invess de todos aqueles que sao do broker especifico
		go func(routeChan string) { errch <- s.channelReadMessageRoutine(newCtx, routeChan) }(channel) // separates several trhead for each channel of this broker
	}

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}

}
