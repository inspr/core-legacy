package kafka

import (
	"context"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// Client is a client for channel operations on kafka
type Client struct {
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
func NewOperator() (operators.ChannelOperatorInterface, error) {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": getEnv().kafkaBootstrapServers,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		k: adminClient,
	}, err
}

// Get gets a channel from kafka
func (c *Client) Get(ctx context.Context, context string, name string) (*meta.Channel, error) {
	topic := toTopic(context, name)
	meta, err := c.k.GetMetadata(&topic, false, 1000)
	if err != nil {
		return nil, ierrors.NewError().InnerError(err).InternalServer().Message("unable to get topic from kafka").Build()
	}

	return fromTopic(name, meta), err

}

// GetAll gets all channels from kafka
func (c *Client) GetAll(ctx context.Context, context string) (ret []*meta.Channel, err error) {
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
func (c *Client) Create(ctx context.Context, context string, channel *meta.Channel) error {
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
func (c *Client) Update(ctx context.Context, context string, channel *meta.Channel) error {
	// updating and creating a new topic is the same thing on kafka
	return c.Create(ctx, context, channel)
}

// Delete deletes a channel from kafka
func (c *Client) Delete(ctx context.Context, context string, name string) error {
	topics := []string{toTopic(context, name)}
	_, err := c.k.DeleteTopics(ctx, topics)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).Message("unable to delete kafka topic").Build()
	}
	return nil
}
