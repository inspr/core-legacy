package sidecarserv

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger = zap.NewNop()

	// logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "loadbalencer-sidecar")))
}

// writeMessageHandler handles requests sent to the write message server
func (s *Server) writeMessageHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling message write")

		channel := strings.TrimPrefix(r.URL.Path, "/")

		if !environment.OutputChannnelList().Contains(channel) {
			logger.Error(fmt.Sprintf("channel %s not found in output channel list", channel))
			insprError := ierrors.NewError().
				BadRequest().
				Message("channel '%s' not found", channel)

			rest.ERROR(w, insprError.Build())
			return
		}

		channelBroker, err := getChannelBroker(channel)
		if err != nil {
			logger.Error("unable to get channel broker",
				zap.String("channel", channel),
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		sidecarAddress := environment.GetBrokerSpecificSidecarAddr(strings.ToUpper(channelBroker))
		sidecarWritePort := environment.GetBrokerWritePort(strings.ToUpper(channelBroker))

		reqAddress := fmt.Sprintf("%s:%s/%s", sidecarAddress, sidecarWritePort, channel)

		logger.Info("sending message to broker",
			zap.String("broker", channelBroker),
			zap.String("channel", channel))

		resp, err := sendRequest(reqAddress, r.Body)
		if err != nil {
			logger.Error("unable to send request to "+channelBroker+" sidecar",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, resp.StatusCode, resp.Body)
	}
}

// readMessageHandler handles requests sent to the read message server
func (s *Server) readMessageHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling message read")

		channel := strings.TrimPrefix(r.URL.Path, "/")

		if !environment.InputChannelList().Contains(channel) {
			logger.Error("channel " + channel + " not found in input channel list")
			insprError := ierrors.NewError().
				BadRequest().
				Message("channel '%s' not found", channel)

			rest.ERROR(w, insprError.Build())
			return
		}

		clientReadPort := os.Getenv("INSPR_SCCLIENT_READ_PORT")
		if clientReadPort == "" {
			insprError := ierrors.NewError().
				NotFound().
				Message("[ENV VAR] INSPR_SCCLIENT_READ_PORT not found")

			rest.ERROR(w, insprError.Build())
			return
		}

		reqAddress := fmt.Sprintf("http://localhost:%v/%v", clientReadPort, channel)

		logger.Info("sending message to node from: ",
			zap.String("channel", channel))

		resp, err := sendRequest(reqAddress, r.Body)
		if err != nil {
			logger.Error("unable to send request to from sidecar to node",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, resp.StatusCode, resp.Body)
	}
}

func getChannelBroker(channel string) (string, error) {
	channelBroker := os.Getenv(channel + "_BROKER")
	if channelBroker == "" {
		return "", ierrors.NewError().
			NotFound().
			Message("[ENV VAR] %v_BROKER not found", channel).
			Build()
	}

	return channelBroker, nil
}

func sendRequest(addr string, body io.Reader) (*http.Response, error) {
	req := http.Client{}
	reqInfo, err := http.NewRequest(http.MethodPost, addr, body)
	if err != nil {
		return nil, err
	}

	resp, err := req.Do(reqInfo)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
