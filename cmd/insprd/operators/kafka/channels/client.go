package channels

import (
	"context"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelOperator is a client for channel operations on kafka
type ChannelOperator struct {
	k *kafka.AdminClient
}

type kafkaEnv struct {
	kafkaBootstrapServers string
}

func getEnv() (env kafkaEnv) {
	boot := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	env.kafkaBootstrapServers = boot
	return
}

// NewOperator returns an initialized operator from the environment variables
func NewOperator() (*ChannelOperator, error) {
	var config *kafka.ConfigMap
	if _, exists := os.LookupEnv("DEBUG"); exists {
		config = &kafka.ConfigMap{
			"bootstrap.servers":       getEnv().kafkaBootstrapServers,
			"test.mock.num.brokers=3": "true",
		}
	} else {
		config = &kafka.ConfigMap{
			"bootstrap.servers": getEnv().kafkaBootstrapServers,
		}
	}
	adminClient, err := kafka.NewAdminClient(config)
	if err != nil {
		return nil, err
	}
	return &ChannelOperator{
		k: adminClient,
	}, err
}

// Get gets a channel from kafka
func (c *ChannelOperator) Get(ctx context.Context, context string, name string) (*meta.Channel, error) {
	topic := toTopic(context, name)
	meta, err := c.k.GetMetadata(&topic, false, 1000)
	if err != nil {
		return nil, ierrors.NewError().InnerError(err).InternalServer().Message("unable to get topic from kafka").Build()
	}

	return fromTopic(name, meta), err

}

// GetAll gets all channels from kafka
func (c *ChannelOperator) GetAll(ctx context.Context, context string) (ret []*meta.Channel, err error) {
	metas, err := c.k.GetMetadata(nil, true, 1000)
	if err != nil {
		return nil, ierrors.NewError().InnerError(err).InternalServer().Message("unable to get topics from kafka").Build()
	}
	for _, topic := range metas.Topics {
		ch := fromTopic(topic.Topic, metas)
		ret = append(ret, ch)
	}
	return
}

// Create creates a channel in kafka
func (c *ChannelOperator) Create(ctx context.Context, context string, channel *meta.Channel) error {
	config, err := configFromChannel(channel)
	if err != nil {
		return err
	}
	_, err = c.k.CreateTopics(ctx, []kafka.TopicSpecification{
		{
			Topic:             toTopic(channel.Meta.Name, context),
			NumPartitions:     config.numberOfPartitions,
			ReplicationFactor: config.replicationFactor,
		},
	})
	if err != nil {
		return ierrors.NewError().InnerError(err).InternalServer().Message("unable to create kafka topic").Build()
	}
	return nil
}

// Update updates a channel in kafka
func (c *ChannelOperator) Update(ctx context.Context, context string, channel *meta.Channel) error {
	// updating and creating a new topic is the same thing on kafka
	return c.Create(ctx, context, channel)
}

// Delete deletes a channel from kafka
func (c *ChannelOperator) Delete(ctx context.Context, context string, name string) error {
	topics := []string{toTopic(context, name)}
	_, err := c.k.DeleteTopics(ctx, topics)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).Message("unable to delete kafka topic").Build()
	}
	return nil
}
