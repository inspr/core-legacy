package lbsidecar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/sidecars/models"
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

		if !environment.OutputChannelList().Contains(channel) {
			logger.Error(fmt.Sprintf("channel %s not found in output channel list", channel))
			insprError := ierrors.NewError().
				BadRequest().
				Message("channel '%s' not found", channel)

			rest.ERROR(w, insprError.Build())
			return
		}

		channelBroker, err := environment.GetChannelBroker(channel)
		if err != nil {
			logger.Error("unable to get channel broker",
				zap.String("channel", channel),
				zap.Any("error", err))

			rest.ERROR(w, err)
			return
		}

		sidecarAddress := environment.GetBrokerSpecificSidecarAddr(channelBroker)
		sidecarWritePort := environment.GetBrokerWritePort(channelBroker)

		reqAddress := fmt.Sprintf("%s:%s/%s", sidecarAddress, sidecarWritePort, channel)

		logger.Debug("encoding message to Avro schema")

		encodedMsg, err := encodeToAvro(channel, r.Body)
		if err != nil {
			logger.Error("unable to encode message to Avro schema",
				zap.String("channel", channel),
				zap.Any("error", err))

			rest.ERROR(w, err)
			return
		}

		logger.Info("sending message to broker",
			zap.String("broker", channelBroker),
			zap.String("channel", channel))

		resp, err := sendRequest(reqAddress, encodedMsg)
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

		logger.Debug("decoding message from Avro schema")

		decodedMsg, err := decodeFromAvro(channel, r.Body)
		if err != nil {
			logger.Error("unable to decode message from Avro schema",
				zap.String("channel", channel),
				zap.Any("error", err))

			rest.ERROR(w, err)
			return
		}

		logger.Info("sending message to node through: ",
			zap.String("channel", channel))

		reqAddress := fmt.Sprintf("http://localhost:%v/%v", clientReadPort, channel)

		resp, err := sendRequest(reqAddress, decodedMsg)
		if err != nil {
			logger.Error("unable to send request to from sidecar to node",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, resp.StatusCode, resp.Body)
	}
}

func sendRequest(addr string, body []byte) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, addr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func encodeToAvro(channel string, body io.Reader) ([]byte, error) {
	var receivedMsg models.BrokerMessage
	json.NewDecoder(body).Decode(&receivedMsg)

	resolvedCh, err := getResolvedChannel(channel)
	if err != nil {
		return nil, err
	}

	encodedAvroMsg, err := encode(resolvedCh, receivedMsg.Data)
	if err != nil {
		return nil, err
	}

	return encodedAvroMsg, nil
}

func decodeFromAvro(channel string, body io.Reader) ([]byte, error) {
	receivedMsg, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	resolvedCh, err := getResolvedChannel(channel)
	if err != nil {
		return nil, err
	}

	decodedAvroMsg, err := readMessage(resolvedCh, receivedMsg)
	if err != nil {
		return nil, err
	}

	jsonEncodedMsg, err := json.Marshal(decodedAvroMsg)
	if err != nil {
		return nil, err
	}

	return jsonEncodedMsg, nil
}

func getResolvedChannel(channel string) (string, error) {
	resolvedCh, ok := os.LookupEnv(channel + "_RESOLVED")
	if !ok {
		logger.Error(fmt.Sprintf("couldn't find resolution for channel %s", channel))
		insprError := ierrors.NewError().
			BadRequest().
			Message("resolution for channel '%s' not found", channel)

		return "", insprError.Build()
	}
	return resolvedCh, nil
}
