package lbsidecar

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
	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/rest"
)

type channelMetric struct {
	messagesRead          prometheus.Counter
	messageSendError      prometheus.Counter
	messageReadError      prometheus.Counter
	messagesSent          prometheus.Counter
	readMessageDuration   prometheus.Summary
	writeMessageDuration  prometheus.Summary
	readMessageThroughput prometheus.Summary
	sendMessageThroughput prometheus.Summary
}

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	writeAddr string
	readAddr  string
	metrics   map[string]channelMetric
}

func (s *Server) GetMetric(channel string) channelMetric {
	metric, ok := s.metrics[channel]
	if ok {
		return metric
	}
	resolved, _ := environment.GetResolvedChannel(channel, environment.GetInputChannelsData(), environment.GetOutputChannelsData())
	broker, _ := environment.GetChannelBroker(channel)
	s.metrics[channel] = channelMetric{
		messagesSent: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "inspr",
			Subsystem: "lbsidecar",
			Name:      "message_send",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
		}),
		messageSendError: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "inspr",
			Subsystem: "lbsidecar",
			Name:      "message_send_error",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
		}),

		messagesRead: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "inspr",
			Subsystem: "lbsidecar",
			Name:      "message_read",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
		}),
		messageReadError: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "inspr",
			Subsystem: "lbsidecar",
			Name:      "messages_read_error",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
		}),

		readMessageDuration: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "lbsidecar",
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
			Subsystem: "lbsidecar",
			Name:      "send_message_duration",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
			Objectives: map[float64]float64{},
		}),
		readMessageThroughput: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "lbsidecar",
			Name:      "read_message_throughput",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
			Objectives: map[float64]float64{},
		}),
		sendMessageThroughput: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "lbsidecar",
			Name:      "send_message_throughput",
			ConstLabels: prometheus.Labels{
				"inspr_channel":          channel,
				"inspr_resolved_channel": resolved,
				"broker":                 broker,
			},
			Objectives: map[float64]float64{},
		}),
	}

	return s.metrics[channel]

}

// Init - initializes a new configured server
func Init() *Server {
	s := Server{}

	wAddr, exists := os.LookupEnv("INSPR_LBSIDECAR_WRITE_PORT")
	if !exists {
		panic("[ENV VAR] INSPR_LBSIDECAR_WRITE_PORT not found")
	}
	rAddr, exists := os.LookupEnv("INSPR_LBSIDECAR_READ_PORT")
	if !exists {
		panic("[ENV VAR] INSPR_LBSIDECAR_READ_PORT not found")
	}

	s.writeAddr = fmt.Sprintf(":%s", wAddr)
	s.readAddr = fmt.Sprintf(":%s", rAddr)
	logger = logger.With(zap.String("read-address", rAddr), zap.String("write-address", wAddr))
	s.metrics = make(map[string]channelMetric)

	return &s
}

// Run starts the server on the port given in addr
func (s *Server) Run(ctx context.Context) error {

	errCh := make(chan error)

	admin := http.NewServeMux()
	admin.Handle("/log/level", alevel)
	admin.Handle("/metrics", promhttp.Handler())
	rest.AttachProfiler(admin)
	adminServer := &http.Server{
		Handler: admin,
		Addr:    "0.0.0.0:16000",
	}
	go func() {
		logger.Info("admin server listening at localhos:16000")
		if err := adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			logger.Error("an error occurred in LB Sidecar write server",
				zap.Error(err))
		}
	}()

	mux := http.NewServeMux()

	mux.Handle("/", s.writeMessageHandler().Post().JSON())

	writeServer := &http.Server{
		Handler: mux,
		Addr:    s.writeAddr,
	}
	go func() {
		if err := writeServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			logger.Error("an error occurred in LB Sidecar write server",
				zap.Error(err))
		}
	}()

	readServer := &http.Server{
		Handler: s.readMessageHandler().Post().JSON(),
		Addr:    s.readAddr,
	}
	go func() {
		if err := readServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			logger.Error("an error occurred in LB Sidecar read server",
				zap.Error(err))
		}
	}()

	logger.Info("LB Sidecar listener is up...")

	select {
	case <-ctx.Done():
		gracefulShutdown(writeServer, readServer, adminServer, nil)
		return ctx.Err()
	case errRead := <-errCh:
		gracefulShutdown(writeServer, readServer, adminServer, errRead)
		return errRead
	}
}

func gracefulShutdown(w, r, a *http.Server, err error) {
	logger.Info("gracefully shutting down...")

	ctxShutdown, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*5),
	)
	defer cancel()

	if err != nil {
		logger.Error("an error occurred in LB Sidecar",
			zap.Error(err))
	}

	// has to be the last method called in the shutdown
	if err = w.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("error while shutting down LB Sidecar write server",
			zap.Error(err))
	}

	if err = r.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("error while shutting down LB Sidecar read server",
			zap.Error(err))
	}

	if err = a.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("error while shutting down LB Sidecar admin server",
			zap.Error(err))
	}
}
