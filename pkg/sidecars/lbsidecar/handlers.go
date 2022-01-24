package lbsidecar

import (
	"bytes"
	//"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/rest"
	"inspr.dev/inspr/pkg/sidecars/models"
)

var logger *zap.Logger
var alevel *zap.AtomicLevel

// const maxBrokerRetries = 5

func init() {
	logger, alevel = logs.Logger(zap.Fields(zap.String("section", "load-balancer-sidecar")))
}

// writeMessageHandler handles requests sent to the write message server
func (s *Server) writeMessageHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		channel := strings.TrimPrefix(r.URL.Path, "/channel/")
		logger.Info("handling message write on " + channel)

		if !environment.OutputChannelList().Contains(channel) {
			logger.Error(fmt.Sprintf("channel %s not found in output channel list", channel))

			rest.ERROR(
				w,
				ierrors.New("channel '%s' not found", channel).BadRequest(),
			)
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

		reqAddress := fmt.Sprintf("%s:%s/channel/%s", sidecarAddress, sidecarWritePort, channel)

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
			s.GetChannelMetric(channel).messageSendError.Inc()

			rest.ERROR(w, err)
			return
		}
		defer resp.Body.Close()

		rest.JSON(w, resp.StatusCode, nil)
		s.GetChannelMetric(channel).messagesSent.Inc()
		elapsed := time.Since(start)
		s.GetChannelMetric(channel).writeMessageDuration.Observe(elapsed.Seconds())
	}
}

func (s *Server) sendRouteRequest() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		path := strings.TrimPrefix(r.URL.Path, "/route/")
		pathArgs := strings.Split(path, "/")
		route := pathArgs[0]
		endpoint := pathArgs[1]

		logger.Info("handling route request", zap.String("route", route), zap.String("path", path))
		resolved, err := environment.GetRouteData(route)

		if err != nil {
			s.getRouteSenderMetric(route).routeSendError.Inc()

			logger.Error("unable to send request to route",
				zap.String("route", route),
				zap.Any("error", err))

			rest.ERROR(w, err)
			return
		}

		if !resolved.Endpoints.Contains(endpoint) {

			s.getRouteSenderMetric(route).routeSendError.Inc()

			err = ierrors.New("invalid endpoint: %s", endpoint).BadRequest()
			logger.Error("unable to send request to "+path,
				zap.Any("error", err))

			rest.ERROR(w, err)
			return
		}
		URL := fmt.Sprintf("%s/route/%s", resolved.Address, path)

		logger.Info("redirecting request", zap.String("route", route), zap.Any("URL", URL))
		http.Redirect(w, r, URL, http.StatusPermanentRedirect)

		elapsed := time.Since(start)
		s.getRouteSenderMetric(route).routeSendDuration.Observe(elapsed.Seconds())
	}
}

// readMessageHandler handles requests sent to the read message server
func (s *Server) readMessageHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		channel := strings.TrimPrefix(r.URL.Path, "/channel/")
		logger.Info("handling message read on " + channel)

		if !environment.InputChannelList().Contains(channel) {
			logger.Error("channel " + channel + " not found in input channel list")
			rest.ERROR(
				w,
				ierrors.New("channel '%s' not found", channel).BadRequest(),
			)
			return
		}

		clientReadPort := os.Getenv("INSPR_SCCLIENT_READ_PORT")
		if clientReadPort == "" {
			rest.ERROR(
				w,
				ierrors.New(
					"[ENV VAR] INSPR_SCCLIENT_READ_PORT not found",
				).NotFound(),
			)
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
			zap.String("channel", channel), zap.String("node port", clientReadPort))

		reqAddress := fmt.Sprintf("http://localhost:%v/channel/%v", clientReadPort, channel)

		resp, err := sendRequest(reqAddress, decodedMsg)
		if err != nil {
			logger.Error("unable to send request from lbsidecar to node",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}
		defer resp.Body.Close()

		rest.JSON(w, resp.StatusCode, resp.Body)
		elapsed := time.Since(start)
		s.GetChannelMetric(channel).readMessageDuration.Observe(elapsed.Seconds())
		s.GetChannelMetric(channel).messagesRead.Add(1)
	}
}

// routeReceiveHandler handles any requests received in the "/route" path, for the lbsidecar
func (s *Server) routeReceiveHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Checking the endpoint
		endpoint := strings.TrimPrefix(r.URL.Path, "/route/")

		splitRoute := strings.SplitN(endpoint, "/", 2)
		if len(splitRoute) == 1 {
			endpoint = ""
		} else {
			endpoint = splitRoute[1]
		}

		// port resolution: using the same as readHandler -> clientReadPort
		clientReadPort := os.Getenv("INSPR_SCCLIENT_READ_PORT")
		if clientReadPort == "" {
			s.GetRouteHandlerMetric(splitRoute[0]).routeReadError.Inc()

			rest.ERROR(
				w,
				ierrors.New(
					"[ENV VAR] INSPR_SCCLIENT_READ_PORT not found",
				).NotFound(),
			)
			return
		}

		// Redirect the request
		// localhost:port/route/endpoint
		client := http.DefaultClient

		URL, _ := url.Parse(fmt.Sprintf("http://localhost:%v/route/%v", clientReadPort, endpoint))
		r.URL = URL
		r.RequestURI = ""
		r.Header.Set("X-Forwarded-For", r.RemoteAddr)

		resp, err := client.Do(r)

		// Validate the response
		if err != nil {
			s.GetRouteHandlerMetric(splitRoute[0]).routeReadError.Inc()

			logger.Error("route: unable to send request from lbsidecar to node",
				zap.Any("error", err))

			rest.ERROR(w, err)
			return
		}
		defer resp.Body.Close()

		elapsed := time.Since(start)
		s.GetRouteHandlerMetric(splitRoute[0]).routeHandleDuration.Observe(elapsed.Seconds())

		// Return the response
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func sendRequest(addr string, body []byte) (*http.Response, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, addr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

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

		return "", ierrors.New(
			"resolution for channel '%s' not found", channel,
		).BadRequest()
	}
	return resolvedCh, nil
}

// func (s *Server) readMessageRoutine(ctx context.Context) error {
// 	// s.runningRead = true
// 	// defer func() { s.runningRead = false }()

// 	errch := make(chan error)
// 	newCtx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	for _, channel := range environment.InputBrokerChannels("broker") {
// 		// separates several threads for each channel of this broker
// 		go func(routeChan string) { errch <- s.channelReadMessageRoutine(newCtx, routeChan) }(
// 			channel,
// 		)
// 	}

// 	select {
// 	case err := <-errch:
// 		return err
// 	case <-ctx.Done():
// 		return ctx.Err()
// 	}

// }

// func (s *Server) channelReadMessageRoutine(
// 	ctx context.Context,
// 	channel string,
// ) error {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		default:
// 			start := time.Now()

// 			var err error
// 			var brokerMsg []byte

// 			brokerMsg, err = s.readWithRetry(ctx, channel)
// 			if err != nil {
// 				return err
// 			}

// 			logger.Debug("trying to send request to loadbalancer",
// 				zap.String("channel", channel),
// 				zap.Any("message", brokerMsg))

// 			status, err := s.writeWithRetry(ctx, channel, brokerMsg)
// 			if err != nil || status != http.StatusOK {
// 				return err
// 			}

// 			// s.Reader.Commit(ctx, channel)
// 			elapsed := time.Since(start)
// 			s.GetChannelMetric(channel).readMessageDuration.Observe(elapsed.Seconds())
// 		}
// 	}
// }

// func (s *Server) readWithRetry(
// 	ctx context.Context,
// 	channel string,
// ) (brokerMsg []byte, err error) {
// 	for i := 0; ; i++ {
// 		brokerMsg, err = s.brokerHandlers["broker"].Reader().ReadMessage(ctx, channel)
// 		if err != nil {
// 			if i == maxBrokerRetries {
// 				return
// 			}
// 			continue
// 		}
// 		return
// 	}
// }

// func (s *Server) writeWithRetry(
// 	ctx context.Context,
// 	channel string,
// 	data []byte,
// ) (status int, err error) {
// 	var resp *http.Response
// 	for i := 0; i <= maxBrokerRetries; i++ {
// 		writeAddr := fmt.Sprintf("%s/channel/%s", s.writeAddr, channel)
// 		logger.Debug("writing with retry",
// 			zap.Any("addr", writeAddr),
// 			zap.Any("write conter", i))

// 		resp, err = s.client.Post(
// 			writeAddr,
// 			"application/octet-stream",
// 			bytes.NewBuffer(data),
// 		)
// 		if resp != nil {
// 			defer resp.Body.Close()
// 		}
// 		status = resp.StatusCode
// 		if err == nil && status == http.StatusOK {
// 			return
// 		}
// 		err = rest.UnmarshalERROR(resp.Body)
// 	}

// 	logger.Debug("unable to send message to lbsidecar",
// 		zap.Error(err))
// 	return
// }
