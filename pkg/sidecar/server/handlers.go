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
	"github.com/inspr/inspr/pkg/sidecar/models"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

// handles the /message route in the server
func (s *Server) writeMessageHandler() rest.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling message write")

		body := models.BrokerData{}
		channel := strings.TrimPrefix(r.URL.Path, "/")

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			insprError := ierrors.NewError().BadRequest().Message("couldn't parse body")
			rest.ERROR(w, insprError.Build())
			return
		}

		if !environment.InputChannelList().Contains(channel) {
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
		if err := s.Writer.WriteMessage(channel, body.Message); err != nil {
			insprError := ierrors.NewError().InternalServer().InnerError(err).Message("broker's writeMessage failed")
			rest.ERROR(w, insprError.Build())
			return
		}
		rest.JSON(w, 200, struct{ Status string }{"OK"})
	}
}

const maxBrokerRetries = 5

func (s *Server) readMessageRoutine(ctx context.Context) error {
	errch := make(chan error)
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	if s.Reader != nil {
		for _, channel := range environment.InputChannelList() {
			// this takes the channel as a parameter to create a new variable from the loop variable
			go func(ctx context.Context, channel string) {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						var err error
						var brokerResp models.BrokerData

						for i := 0; ; i++ {
							fmt.Println(channel, ": reading message")
							brokerResp, err = s.Reader.ReadMessage(ctx, channel)
							fmt.Println(channel, ": message read")
							insprError := ierrors.NewError().InternalServer().InnerError(err).Message("broker's ReadMessage returned an error").Build()
							if err != nil {
								if i == maxBrokerRetries {
									select {
									case errch <- insprError:
									case <-ctx.Done():
									}
									return
								}
								logger.Info("error reading message from broker", zap.Any("error", insprError))
								continue
							}
							break
						}

						type response struct {
							Status string
						}
						resp := response{}

						fmt.Println("trying to send requess")

						err = s.client.Send(ctx, "/"+channel, http.MethodPost, brokerResp, &resp)
						if err != nil {
							logger.Info("error sending message to dapp", zap.Any("error", err))
						} else {
							s.Reader.Commit(ctx, channel)
						}
					}
				}
			}(newCtx, channel)
		}
	}
	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		return nil
	}

}

// Close closes the server connection
func (s *Server) Close() {
	s.cancel()
}
