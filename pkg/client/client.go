package dappclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/rest"
	"inspr.dev/inspr/pkg/rest/request"
	"inspr.dev/inspr/pkg/sidecars/models"
)

var logger *zap.Logger
var alevel *zap.AtomicLevel

func init() {
	logger, alevel = logs.Logger(zap.Fields(zap.String("section", "sidecar-client"), zap.String("dapp-name", os.Getenv("INSPR_APP_ID"))))
}

// Client is the struct which implements the methods of AppClient interface
type Client struct {
	client   *request.Client
	mux      *http.ServeMux
	readAddr string
}

// NewAppClient returns a new instance of the client of the AppClient package
func NewAppClient() *Client {
	logger.Info("initializing dapp client")
	writeAddr := fmt.Sprintf("http://localhost:%s", os.Getenv("INSPR_LBSIDECAR_WRITE_PORT"))
	readAddr := fmt.Sprintf(":%s", os.Getenv("INSPR_SCCLIENT_READ_PORT"))
	logger.Info("got configuration from environment variables")
	logger = logger.With(zap.String("read-address", readAddr), zap.String("write-address", writeAddr))
	return &Client{
		readAddr: readAddr,
		client: request.NewClient().
			BaseURL(writeAddr).
			Encoder(json.Marshal).
			Decoder(request.JSONDecoderGenerator).
			Pointer(),
		mux: http.NewServeMux(),
	}
}

// WriteMessage receives a channel and a message and sends it in a request to the sidecar server
func (c *Client) WriteMessage(ctx context.Context, channel string, msg interface{}) error {
	l := logger.With(zap.String("operation", "write"), zap.String("channel", channel))
	l.Info("received write message request")
	data := models.BrokerMessage{
		Data: msg,
	}

	var resp interface{}
	// sends a message to the corresponding channel route on the sidecar
	l.Debug("sending message to load balancer")
	err := c.client.Send(
		ctx,
		"/channel/"+channel,
		http.MethodPost,
		data,
		&resp)
	if err != nil {
		l.Error("error sending message to load balancer")
	} else {
		l.Info("message sent")
	}
	return err
}

// HandleChannel handles messages received in a given channel.
func (c *Client) HandleChannel(channel string, handler func(ctx context.Context, body io.Reader) error) {
	c.addHandlerToMux("/channel/"+channel, handler)
}

// HandleRoute handles messages received in a given route.
func (c *Client) HandleRoute(path string, handler func(ctx context.Context, body io.Reader) error) {
	path = strings.TrimPrefix(path, "/")
	c.addHandlerToMux("/route/"+path, handler)
}

//Run runs the server with the handlers defined in HandleChannel
func (c *Client) Run(ctx context.Context) error {

	var err error
	c.mux.Handle("/log/level", alevel)
	server := http.Server{
		Handler: c.mux,
		Addr:    c.readAddr,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("error serving dApp", zap.Error(err))
		}
	}()

	logger.Info("inspr client server is running", zap.String("log-level", alevel.String()))

	<-ctx.Done()

	logger.Info("gracefully shutting down")

	ctxShutdown, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*5),
	)
	defer cancel()

	if err != nil {
		logger.Fatal("error in server shitting down", zap.Error(err))
	}

	// has to be the last method called in the shutdown
	if err = server.Shutdown(ctxShutdown); err != nil {
		return err
	}
	return ctx.Err()
}

func (c *Client) addHandlerToMux(path string, handler func(ctx context.Context, body io.Reader) error) {
	c.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received request in client handler", zap.String("Path", path))
		// user defined handler. Returns error if the user wants to return it
		err := handler(context.Background(), r.Body)
		if err != nil {
			logger.Error("error returned by client handler", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		rest.JSON(w, 200, nil)
	})
}
