package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"inspr.dev/inspr/examples/kafka_standalone/vars"
)

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": vars.BootstrapServers,
		"group.id":          "foo",
		"auto.offset.reset": "earliest"})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err.Error())
		return
	}

	err = consumer.Subscribe(vars.Topic, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe consumer: %s\n", err.Error())
		return
	}

	readMsgDuration := promauto.NewSummary(prometheus.SummaryOpts{
		Namespace: "kafkasa",
		Subsystem: "reader",
		Name:      "read_message_duration",
		ConstLabels: prometheus.Labels{
			"broker": "kafka",
		},
		Objectives: map[float64]float64{},
	})

	admin := http.NewServeMux()
	admin.Handle("/metrics", promhttp.Handler())
	adminServer := &http.Server{
		Handler: admin,
		Addr:    "0.0.0.0:16000",
	}

	run := true

	go func() {
		fmt.Println("admin server listening at localhost:16000")
		if err := adminServer.ListenAndServe(); err != nil {
			run = false
			fmt.Println("an error occurred in LB Sidecar write server", err)
		}
	}()

	start := time.Now()
	for run {
		ev := consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("Message on %s: %s\n",
				vars.Topic, string(e.Value))
			elapsed := time.Since(start)
			readMsgDuration.Observe(elapsed.Seconds())
			start = time.Now()
		case kafka.PartitionEOF:
			fmt.Printf("Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "Error: %v\n", e)
			run = false
		default:
		}
	}

	consumer.Close()
}
