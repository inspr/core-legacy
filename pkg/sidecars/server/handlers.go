package sidecarserv

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"go.uber.org/zap"
)

const maxBrokerRetries = 5

// handles the messages route in the server
func (s *Server) writeMessageHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling message write")

		channel := strings.TrimPrefix(r.URL.Path, "/")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		if !environment.OutputBrokerChannnels(s.broker).Contains(channel) {
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
			rest.ERROR(w, ierrors.NewError().Message("broker's writeMessage failed, %s", err.Error()).Build())
			return
		}
		rest.JSON(w, 200, nil)
	}
}

func (s *Server) readMessageRoutine(ctx context.Context) error {
	s.runningRead = true
	defer func() { s.runningRead = false }()

	errch := make(chan error)
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, channel := range environment.InputBrokerChannels(s.broker) {
		// separates several trhead for each channel of this broker
		go func(routeChan string) { errch <- s.channelReadMessageRoutine(newCtx, routeChan) }(channel)
	}

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}

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

			status, err := s.writeWithRetry(ctx, channel, brokerMsg)
			if err != nil || status != http.StatusOK {
				return err
			}
			s.Reader.Commit(ctx, channel)
		}
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

func (s *Server) writeWithRetry(ctx context.Context, channel string, data []byte) (status int, err error) {
	var resp *http.Response
	for i := 0; i <= maxBrokerRetries; i++ {
		resp, err = s.client.Post(s.outAddr, "application/octet-stream", bytes.NewBuffer(data))
		status = resp.StatusCode
		if err == nil && status == http.StatusOK {
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&status)
			return
		}
		err = rest.UnmarshalERROR(resp.Body)
	}
	return
}
