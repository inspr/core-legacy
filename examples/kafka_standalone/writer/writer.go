package main

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()
	time.Sleep(5 * time.Second)

	// delivery_chan := make(chan kafka.Event, 10000)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &vars.Topic,
			Partition: kafka.PartitionAny},
		Value: []byte("hello_world")},
		nil,
	)
	if err != nil {
		fmt.Printf("Failed to produce message: %s\n", err.Error())
		return
	}
}
