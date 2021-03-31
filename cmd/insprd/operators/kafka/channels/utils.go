package channels

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type kafkaConfiguration struct {
	numberOfPartitions int
	replicationFactor  int
}

func configFromChannel(ch *meta.Channel) (kafkaConfiguration, error) {

	config := kafkaConfiguration{}
	if nPart, ok := ch.Meta.Annotations["kafka.partition.number"]; ok {
		var err error
		config.numberOfPartitions, err = strconv.Atoi(nPart)
		if err != nil {
			config.numberOfPartitions = 1
			return config, ierrors.NewError().
				InvalidChannel().
				Message(
					"invalid partition configuration %s",
					ch.Meta.Annotations["kafka.partition.number"],
				).
				Build()
		}
	}

	if nPart, ok := ch.Meta.Annotations["kafka.replication.factor"]; ok {
		var err error
		config.replicationFactor, err = strconv.Atoi(nPart)
		if err != nil {
			config.replicationFactor = 1
			return config, ierrors.NewError().
				InvalidChannel().
				Message(
					"invalid replication configuration %s",
					ch.Meta.Annotations["kafka.replication.factor"],
				).
				Build()
		}
	}

	return config, nil
}

func toTopic(ctx, name string) string {
	insprEnvironment := environment.GetInsprEnvironment()
	if insprEnvironment == "" {
		return fmt.Sprintf("inspr-%s-%s", ctx, name)
	}
	return fmt.Sprintf("inspr-%s-%s-%s", insprEnvironment, ctx, name)
}

func fromTopic(name string, meta *kafka.Metadata) (ch *meta.Channel) {
	ch.Meta.Annotations["kafka.partition.number"] = strconv.Itoa(len(meta.Topics[name].Partitions))
	splitName := strings.Split(name, "-")
	if len(splitName) == 4 {
		ch.Meta.Name = splitName[3]
		ch.Meta.Parent = splitName[2]
	} else if len(splitName) == 3 {
		ch.Meta.Name = splitName[2]
		ch.Meta.Parent = splitName[1]
	}
	return
}
