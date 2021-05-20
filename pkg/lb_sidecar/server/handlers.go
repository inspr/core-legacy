package sidecarserv

import (
	"context"
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

// handles requests sent to the write message server
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

		resp, err := sendWriteRequest(reqAddress, r.Body)
		if err != nil {
			logger.Error("unable to send request to "+channelBroker+" sidecar",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, resp.StatusCode, resp.Body)
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

// Write message helper functions

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

func sendWriteRequest(addr string, body io.Reader) (*http.Response, error) {
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
