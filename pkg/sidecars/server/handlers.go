package sidecarserv

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/rest"
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
			rest.ERROR(
				w,
				ierrors.New("channel '%s' not found", channel).BadRequest(),
			)
			return
		}

		logger.Info("writing message to broker",
			zap.String("channel", channel))
		if err := s.Writer.WriteMessage(channel, body); err != nil {
			rest.ERROR(
				w,
				ierrors.New("broker's WriteMessage failed, %s", err.Error()),
			)
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
		// separates several threads for each channel of this broker
		go func(routeChan string) { errch <- s.channelReadMessageRoutine(newCtx, routeChan) }(
			channel,
		)
	}

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}

}

func (s *Server) channelReadMessageRoutine(
	ctx context.Context,
	channel string,
) error {
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

			logger.Debug("trying to send request to loadbalancer",
				zap.String("channel", channel),
				zap.Any("message", brokerMsg))

			status, err := s.writeWithRetry(ctx, channel, brokerMsg)
			if err != nil || status != http.StatusOK {
				return err
			}
			s.Reader.Commit(ctx, channel)
		}
	}
}

func (s *Server) readWithRetry(
	ctx context.Context,
	channel string,
) (brokerMsg []byte, err error) {
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

func (s *Server) writeWithRetry(
	ctx context.Context,
	channel string,
	data []byte,
) (status int, err error) {
	var resp *http.Response
	for i := 0; i <= maxBrokerRetries; i++ {
		writeAddr := fmt.Sprintf("%s/%s", s.outAddr, channel)
		logger.Debug("writing with retry",
			zap.Any("addr", writeAddr),
			zap.Any("write conter", i))

		resp, err = s.client.Post(
			writeAddr,
			"application/octet-stream",
			bytes.NewBuffer(data),
		)
		if resp != nil {
			defer resp.Body.Close()
		}

		status = resp.StatusCode
		if err == nil && status == http.StatusOK {
			return
		}
		defer resp.Body.Close()
		err = rest.UnmarshalERROR(resp.Body)
	}

	logger.Debug("unable to send message to lbsidecar",
		zap.Error(err))
	return
}
