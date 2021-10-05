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
	messagesRead         prometheus.Counter
	messageSendError     prometheus.Counter
	messageReadError     prometheus.Counter
	messagesSent         prometheus.Counter
	readMessageDuration  prometheus.Summary
	writeMessageDuration prometheus.Summary
}

type routeMetric struct {
	routeReadError      prometheus.Counter
	routeSendError      prometheus.Counter
	routesendDuration   prometheus.Summary
	routeHandleDuration prometheus.Summary
}

// Server is a struct that contains the variables necessary
// to handle the necessary routes of the rest API
type Server struct {
	writeAddr     string
	readAddr      string
	channelMetric map[string]channelMetric
	routeMetric   map[string]routeMetric
}

func (s *Server) GetMetricRoute(route string) routeMetric {
	metric, ok := s.routeMetric[route]
	if ok {
		return metric
	}

	resolved, _ := environment.GetRouteData(route)
	s.routeMetric[route] = routeMetric{
		routeSendError: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "inspr",
			Subsystem: "lbsidecar",
			Name:      "message_send_error",
			ConstLabels: prometheus.Labels{
				"inspr_route":          route,
				"inspr_route_adress":   resolved.Address,
				"inspr_route_endpoint": resolved.Endpoints.Join(";"),
			},
		}),
		routeReadError: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "inspr",
			Subsystem: "lbsidecar",
			Name:      "message_read_error",
			ConstLabels: prometheus.Labels{
				"inspr_route":          route,
				"inspr_route_adress":   resolved.Address,
				"inspr_route_endpoint": resolved.Endpoints.Join(";"),
			},
		}),

		routesendDuration: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "lbsidecar",
			Name:      "route_send_duration",
			ConstLabels: prometheus.Labels{
				"inspr_route":          route,
				"inspr_route_adress":   resolved.Address,
				"inspr_route_endpoint": resolved.Endpoints.Join(";"),
			},
			Objectives: map[float64]float64{},
		}),
		routeHandleDuration: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: "inspr",
			Subsystem: "lbsidecar",
			Name:      "route_handle_duration",
			ConstLabels: prometheus.Labels{
				"inspr_route":          route,
				"inspr_route_adress":   resolved.Address,
				"inspr_route_endpoint": resolved.Endpoints.Join(";"),
			},
			Objectives: map[float64]float64{},
		}),
	}

	return s.routeMetric[route]

}

func (s *Server) GetMetricChannel(channel string) channelMetric {
	metric, ok := s.channelMetric[channel]
	if ok {
		return metric
	}

	resolved, _ := environment.GetResolvedChannel(channel, environment.GetInputChannelsData(), environment.GetOutputChannelsData())
	broker, _ := environment.GetChannelBroker(channel)
	s.channelMetric[channel] = channelMetric{
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
	}

	return s.channelMetric[channel]

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
	s.channelMetric = make(map[string]channelMetric)
	s.routeMetric = make(map[string]routeMetric)

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
		logger.Info("admin server listening at localhost:16000")
		if err := adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			logger.Error("an error occurred in LB Sidecar write server",
				zap.Error(err))
		}
	}()

	muxWriter := http.NewServeMux()

	muxWriter.Handle("/channel/", s.writeMessageHandler().Post().JSON())
	muxWriter.Handle("/route/", s.sendRequest().JSON())

	writeServer := &http.Server{
		Handler: muxWriter,
		Addr:    s.writeAddr,
	}
	go func() {
		if err := writeServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			logger.Error("an error occurred in LB Sidecar write server",
				zap.Error(err))
		}
	}()

	muxReader := http.NewServeMux()

	muxReader.Handle("/channel/", s.readMessageHandler().Post().JSON())
	muxReader.Handle("/route/", s.routeReceiveHandler().JSON())

	readServer := &http.Server{
		Handler: muxReader,
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
