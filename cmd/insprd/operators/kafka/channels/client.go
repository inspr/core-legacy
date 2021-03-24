package channels

import (
	"context"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"go.uber.org/zap"
)

// ChannelOperator is a client for channel operations on kafka
type ChannelOperator struct {
	k      *kafka.AdminClient
	logger *zap.Logger
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
	logConf := zap.NewDevelopmentConfig()
	logger, _ := logConf.Build(zap.Fields(zap.String("section", "kafka-operator")))

	var config *kafka.ConfigMap

	if _, exists := os.LookupEnv("DEBUG"); exists {
		logger.Info("initializing kafka admin with debug configs")
		config = &kafka.ConfigMap{
			"bootstrap.servers":     getEnv().kafkaBootstrapServers,
			"test.mock.num.brokers": "3",
		}
	} else {
		logger.Info("initializing kafka admin with production configs", zap.String("kafka bootstrap servers", getEnv().kafkaBootstrapServers))
		config = &kafka.ConfigMap{
			"bootstrap.servers": getEnv().kafkaBootstrapServers,
		}
	}

	adminClient, err := kafka.NewAdminClient(config)
	if err != nil {
		logger.Error("unable to create kafka admin client", zap.Any("error", err))
		return nil, err
	}
	return &ChannelOperator{
		k:      adminClient,
		logger: logger,
	}, err
}

// Get gets a channel from kafka
func (c *ChannelOperator) Get(ctx context.Context, context string, name string) (*meta.Channel, error) {
	topic := toTopic(context, name)
	meta, err := c.k.GetMetadata(&topic, false, 1000)
	if err != nil {
		log.Println(err)
		return nil, ierrors.NewError().InnerError(err).InternalServer().Message("unable to get topic from kafka").Build()
	}

	return fromTopic(name, meta), err

}

// GetAll gets all channels from kafka
func (c *ChannelOperator) GetAll(ctx context.Context, context string) (ret []*meta.Channel, err error) {
	metas, err := c.k.GetMetadata(nil, true, 1000)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return err
	}
	configs := []kafka.TopicSpecification{
		{
			Topic:             toTopic(channel.Meta.Name, context),
			NumPartitions:     config.numberOfPartitions,
			ReplicationFactor: config.replicationFactor,
		},
	}
	_, err = c.k.CreateTopics(ctx, configs)
	if err != nil {
		c.logger.Error("error creating kafka topic", zap.Any("error", err))
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
		log.Println(err)
		return ierrors.NewError().InternalServer().InnerError(err).Message("unable to delete kafka topic").Build()
	}
	return nil
}
