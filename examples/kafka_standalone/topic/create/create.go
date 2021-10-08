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

	configs := []kafka.TopicSpecification{
		{
			Topic:             vars.Topic,
			NumPartitions:     vars.NumPartitions,
			ReplicationFactor: vars.ReplicationFactor,
		},
	}
	_, err = adminClient.CreateTopics(context.Background(), configs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
