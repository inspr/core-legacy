package channels

import (
	"context"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"go.uber.org/zap"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "kafka-channel-operator")))
}

// ChannelOperator is a client for channel operations on kafka
type ChannelOperator struct {
	k      kafkaAdminClient
	logger *zap.Logger
	mem    memory.Manager
}

type kafkaEnv struct {
	kafkaBootstrapServers string
}

func getEnv() (env kafkaEnv) {
	boot := os.Getenv("INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS")
	env.kafkaBootstrapServers = boot
	return
}

// NewOperator returns an initialized operator from the environment variables
func NewOperator(mem memory.Manager) (*ChannelOperator, error) {

	var config *kafka.ConfigMap
	var err error
	var adminClient kafkaAdminClient
	if _, exists := os.LookupEnv("DEBUG"); exists {
		logger.Info("initializing kafka admin with debug configs")
		adminClient = &mockAdminClient{}
	} else {
		logger.Info("initializing kafka admin with production configs",
			zap.String("kafka bootstrap servers", "kafka.default.svc:9092"))
		config = &kafka.ConfigMap{
			"bootstrap.servers": "kafka.default.svc:9092",
		}

		adminClient, err = kafka.NewAdminClient(config)
		if err != nil {
			logger.Error("unable to create kafka admin client", zap.Any("error", err))
			return nil, err
		}
	}

	return &ChannelOperator{
		k:      adminClient,
		logger: logger,
		mem:    mem,
	}, err
}

// Get gets a channel from kafka
func (c *ChannelOperator) Get(ctx context.Context, context string, name string) (*meta.Channel, error) {
	channel, _ := c.mem.Root().Channels().Get(context, name)
	logger.Info("trying to get Channel from Kafka Topic",
		zap.String("channel", name),
		zap.String("context", context))

	topic := toTopic(channel)
	meta, err := c.k.GetMetadata(&topic, false, 1000)
	if err != nil {
		logger.Error("unable to get Kafka Topic", zap.Any("error", err))
		return nil, ierrors.NewError().InnerError(err).InternalServer().Message("unable to get topic from kafka").Build()
	}

	return fromTopic(name, meta), err

}

// GetAll gets all channels from kafka
func (c *ChannelOperator) GetAll(ctx context.Context, context string) (ret []*meta.Channel, err error) {
	logger.Info("trying to get all Channels from Kafka Topics",
		zap.String("context", context))

	metas, err := c.k.GetMetadata(nil, true, 1000)
	if err != nil {
		logger.Error("unable to get all Kafka Topics", zap.Any("error", err))
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
	logger.Info("trying to create a Channel in Kafka",
		zap.String("channel", channel.Meta.Name),
		zap.String("context", context))

	config, err := configFromChannel(channel)
	if err != nil {
		logger.Error("unable to extract Kafka config from given Channel",
			zap.Any("error", err))
		return err
	}

	configs := []kafka.TopicSpecification{
		{
			Topic:             toTopic(channel),
			NumPartitions:     config.numberOfPartitions,
			ReplicationFactor: config.replicationFactor,
		},
	}
	_, err = c.k.CreateTopics(ctx, configs)
	if err != nil {
		logger.Error("error creating Kafka Topic", zap.Any("error", err))
		return ierrors.NewError().InnerError(err).InternalServer().Message("unable to create kafka topic").Build()
	}
	return nil
}

// Update updates a channel in kafka
func (c *ChannelOperator) Update(ctx context.Context, context string, channel *meta.Channel) error {
	logger.Info("trying to update a Channels in Kafka",
		zap.String("channel", channel.Meta.Name),
		zap.String("context", context))
	// updating and creating a new topic is the same thing on kafka
	return c.Create(ctx, context, channel)
}

// Delete deletes a channel from kafka
func (c *ChannelOperator) Delete(ctx context.Context, context string, name string) error {
	channel, _ := c.mem.Root().Channels().Get(context, name)
	topics := []string{toTopic(channel)}
	logger.Info("trying to delete a Channel from Kafka Topics",
		zap.String("channel", name),
		zap.String("context", context))

	_, err := c.k.DeleteTopics(ctx, topics)
	if err != nil {
		logger.Error("error deleting Kafka Topic", zap.Any("error", err))
		return ierrors.NewError().InternalServer().InnerError(err).Message("unable to delete kafka topic").Build()
	}
	return nil
}

type kafkaAdminClient interface {
	DeleteTopics(ctx context.Context, topics []string, options ...kafka.DeleteTopicsAdminOption) (result []kafka.TopicResult, err error)
	CreateTopics(ctx context.Context, topics []kafka.TopicSpecification, options ...kafka.CreateTopicsAdminOption) (result []kafka.TopicResult, err error)
	GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error)
}

type mockAdminClient struct {
}

func (*mockAdminClient) DeleteTopics(ctx context.Context, topics []string, options ...kafka.DeleteTopicsAdminOption) (result []kafka.TopicResult, err error) {
	return nil, nil
}

func (*mockAdminClient) CreateTopics(ctx context.Context, topics []kafka.TopicSpecification, options ...kafka.CreateTopicsAdminOption) (result []kafka.TopicResult, err error) {
	return nil, nil
}
func (*mockAdminClient) GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error) {
	return nil, nil
}
