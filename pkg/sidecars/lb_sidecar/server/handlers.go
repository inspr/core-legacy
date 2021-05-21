package sidecarserv

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/inspr/inspr/cmd/insprd/memory/tree"
	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/ierrors"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/rest"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "loadbalencer-sidecar")))
}

// writeMessageHandler handles requests sent to the write message server
func (s *Server) writeMessageHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling message write")

		channel := strings.TrimPrefix(r.URL.Path, "/")

		if !environment.OutputChannnelList().Contains(channel) {
			logger.Error("channel " + channel + " not found in output channel list")
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

		sidecarWritePort := os.Getenv("INSPR_SIDECAR_" + strings.ToUpper(channelBroker) + "_WRITE_PORT")
		if sidecarWritePort == "" {
			logger.Error("unable to get broker " + channelBroker + " port")
			insprError := ierrors.NewError().
				NotFound().
				Message("[ENV VAR] INSPR_SIDECAR_%s_WRITE_PORT not found", channelBroker)

			rest.ERROR(w, insprError.Build())
			return
		}

		reqAddress := fmt.Sprintf("http://localhost:%v/%v", sidecarWritePort, channel)

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
	pathToChannel := os.Getenv(channel + "_RESOLVED_SCOPE")
	if pathToChannel != "" {
		return "", ierrors.NewError().
			NotFound().
			Message("[ENV VAR] %s_RESOLVED_SCOPE not found", channel).
			Build()
	}

	scope, chName, err := metautils.RemoveLastPartInScope(pathToChannel)
	if err != nil {
		return "", err
	}

	chStructure, err := tree.GetTreeMemory().Channels().Get(scope, chName)
	if err != nil {
		return "", err
	}

	return chStructure.Spec.SelectedBroker, nil
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
