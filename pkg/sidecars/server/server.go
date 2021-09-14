package sidecarserv

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/rest"
	"inspr.dev/inspr/pkg/sidecars/models"
)

type channelMetric struct {
	readMessageDuration  prometheus.Summary
	writeMessageDuration prometheus.Summary
	readTimeDuration     prometheus.Summary
	writeTimeDuration    prometheus.Summary
}

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	broker       string
	Reader       models.Reader
	Writer       models.Writer
	inAddr       string
	outAddr      string
	client       *http.Client
	runningRead  bool
	runningWrite bool
	metrics      map[string]channelMetric
}

var logger *zap.Logger
var alevel *zap.AtomicLevel

func init() {
	logger, alevel = logs.Logger(zap.Fields(zap.String("section", "sidecar")))
}

// Init - configures the server
func Init(r models.Reader, w models.Writer, broker string) *Server {
	server := &Server{}
	// server fetches addresses variable names from models.
	envVars := brokers.GetSidecarConnectionVars(broker)
	if envVars == nil {
		panic(fmt.Sprintf("%s broker's enviroment variables not configured", broker))
	}
	server.broker = broker

	// server fetches required addresses from deployment.
	inAddr, ok := os.LookupEnv(envVars.WriteEnvVar)
	if !ok {
		panic(fmt.Sprintf("[ENV VAR] %s not found", envVars.WriteEnvVar))
	}

	outAddr, ok := os.LookupEnv(envVars.ReadEnvVar)
	if !ok {
		panic(fmt.Sprintf("[ENV VAR] %s not found", envVars.ReadEnvVar))
	}

	server.inAddr = fmt.Sprintf(":%s", inAddr)
	server.outAddr = fmt.Sprintf("http://localhost:%v", outAddr)
	server.client = &http.Client{}

	// implementations of write and read for a specific sidecar
	server.Reader = r
	server.Writer = w
	server.metrics = make(map[string]channelMetric)
	return server
}

func (s *Server) GetMetric(channel string) channelMetric {
	metric, ok := s.metrics[channel]
	if ok {
		return metric
	}
	resolved, _ := environment.GetResolvedChannel(channel, environment.GetInputChannelsData(), environment.GetOutputChannelsData())
	broker, _ := environment.GetChannelBroker(channel)
	s.metrics[channel] = channelMetric{
		readMessageDuration: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "sidecar",
			Name:      "read_message_duration",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
			Objectives: map[float64]float64{},
		}),
		writeMessageDuration: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "sidecar",
			Name:      "send_message_duration",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
			Objectives: map[float64]float64{},
		}),

		readTimeDuration: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "sidecar",
			Name:      "read_message_duration",
			ConstLabels: prometheus.Labels{
				"inspr_app_id":           environment.GetInsprAppID(),
				"isnpr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
			Objectives: make(map[float64]float64),
		}),

		writeTimeDuration: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "sidecar",
			Name:      "read_message_duration",
			ConstLabels: prometheus.Labels{
				"inspr_app_id":           environment.GetInsprAppID(),
				"isnpr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
			Objectives: make(map[float64]float64),
		}),
	}

	return s.metrics[channel]

}

// Run starts the server on the port given in addr
func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.Handle("/", s.writeMessageHandler().Post().JSON())

	server := &http.Server{
		Handler: mux,
		Addr:    s.inAddr,
	}

	errCh := make(chan error)
	admin := http.NewServeMux()
	admin.Handle("/log/level", alevel)
	admin.Handle("/metrics", promhttp.Handler())
	rest.AttachProfiler(admin)
	adminServer := &http.Server{
		Handler: admin,
		Addr:    "0.0.0.0:16001",
	}
	go func() {
		logger.Info("admin server listening at localhos:16001")
		if err := adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			logger.Error("an error occurred in LB Sidecar write server",
				zap.Error(err))
		}
	}()
	// create read message routine and captures its error
	go func() { errCh <- s.readMessageRoutine(ctx) }()

	var err error
	go func() {
		s.runningWrite = true
		defer func() { s.runningWrite = false }()
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("an error ocurred in %v sidecar: %v", s.broker, err))
			errCh <- err
		}
	}()

	logger.Info(fmt.Sprintf("%s sidecar listener is up...", s.broker))

	select {
	case <-ctx.Done():
		s.gracefulShutdown(server, ctx.Err())
		return ctx.Err()
	case errRead := <-errCh:
		s.gracefulShutdown(server, errRead)
		return errRead
	}

}

func (s *Server) gracefulShutdown(server *http.Server, err error) {
	logger.Info("gracefully shutting down...")

	ctxShutdown, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*5),
	)

	defer cancel()

	if err != nil {
		logger.Error("an error occurred on sidecar",
			zap.Any("broker", s.broker), zap.Error(err))
	}

	s.Writer.Close()

	// has to be the last method called in the shutdown
	if err = server.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("error shutting down server")
	}
}
