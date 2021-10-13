package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"inspr.dev/inspr/examples/kafka_standalone/vars"
)

func main() {

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": vars.BootstrapServers,
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err.Error())
		return
	}

	writeMsgDuration := promauto.NewSummary(prometheus.SummaryOpts{
		Namespace: "kafkasa",
		Subsystem: "writer",
		Name:      "write_message_duration",
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
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
					run = false
				} else {
					fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
					elapsed := time.Since(start)
					writeMsgDuration.Observe(elapsed.Seconds())
					start = time.Now()
				}
			}
		}
	}()

	for run {
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &vars.Topic,
				Partition: kafka.PartitionAny},
			Value: []byte("hello_world")},
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to produce message: %s\n", err.Error())
			run = false
			return
		}
	}
}
