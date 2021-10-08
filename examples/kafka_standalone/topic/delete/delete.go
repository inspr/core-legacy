package main

import (
	"context"
	"fmt"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"inspr.dev/inspr/examples/kafka_standalone/vars"
)

func main() {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": vars.BootstrapServers,
	}

	adminClient, err := kafka.NewAdminClient(kafkaConfig)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = adminClient.DeleteTopics(context.Background(), []string{vars.Topic})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
