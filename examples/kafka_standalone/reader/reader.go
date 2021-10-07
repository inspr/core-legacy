package main

import (
	"fmt"
	"os"

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

	consumer.Subscribe(vars.Topic, nil)

	run := true
	for run {
		ev := consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("%% Message on %s:\n%s\n",
				e.TopicPartition, string(e.Value))
		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}

	consumer.Close()
}
