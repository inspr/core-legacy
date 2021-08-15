package kafkaop

import (
	"context"
	"os"

	"go.uber.org/zap"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
	"inspr.dev/inspr/cmd/sidecars"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/meta"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = logs.Logger(
		zap.Fields(zap.String("section", "kafka-channel-operator")),
	)
}

// ChannelOperator is a client for channel operations on kafka
type ChannelOperator struct {
	k      kafkaAdminClient
	logger *zap.Logger
	mem    tree.Manager
}

// NewOperator returns an initialized operator from the environment variables
func NewOperator(
	mem tree.Manager,
	config sidecars.KafkaConfig,
) (*ChannelOperator, error) {
	logger.Debug("initializing operator")
	var kafkaConfig *kafka.ConfigMap
	var err error
	var adminClient kafkaAdminClient
	if _, exists := os.LookupEnv("DEBUG"); exists {
		logger.Debug("initializing kafka admin with debug configs")
		adminClient = &mockAdminClient{}
	} else {
		bootstrap := config.BootstrapServers
		logger.Debug("initializing kafka admin with production configs",
			zap.String("kafka-bootstrap-servers", bootstrap))
		kafkaConfig = &kafka.ConfigMap{
			"bootstrap.servers": bootstrap,
		}

		adminClient, err = kafka.NewAdminClient(kafkaConfig)
		if err != nil {
			logger.Error("unable to create kafka admin client", zap.Error(err))
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
func (c *ChannelOperator) Get(
	ctx context.Context,
	context string,
	name string,
) (*meta.Channel, error) {
	channel, _ := c.mem.Perm().Channels().Get(context, name)
	l := logger.With(
		zap.String("channel", name),
		zap.String("context", context))

	l.Debug("trying to get Channel from Kafka Topic")

	topic := toTopic(channel)
	meta, err := c.k.GetMetadata(&topic, false, 1000)
	if err != nil {
		l.Error("unable to get Kafka Topic", zap.Error(err))
		return nil, ierrors.Wrap(
			ierrors.New(err).InternalServer(),
			"unable to get topic from kafka",
		)
	}

	return fromTopic(name, meta), err

}

// GetAll gets all channels from kafka
func (c *ChannelOperator) GetAll(
	ctx context.Context,
	context string,
) (ret []*meta.Channel, err error) {
	logger.Info("trying to get all Channels from Kafka Topics",
		zap.String("context", context))

	metas, err := c.k.GetMetadata(nil, true, 1000)
	if err != nil {
		logger.Error("unable to get all Kafka Topics", zap.Any("error", err))
		return nil, ierrors.Wrap(
			ierrors.New(err).InternalServer(),
			"unable to get topics from kafka",
		)
	}
	for _, topic := range metas.Topics {
		ch := fromTopic(topic.Topic, metas)
		ret = append(ret, ch)
	}
	return
}

// Create creates a channel in kafka
func (c *ChannelOperator) Create(
	ctx context.Context,
	context string,
	channel *meta.Channel,
) error {
	l := logger.With(
		zap.String("channel", channel.Meta.Name),
		zap.String("context", context),
	)
	l.Info("trying to create a Channel in Kafka")

	config, err := configFromChannel(channel)
	if err != nil {
		l.Error(
			"unable to extract Kafka config from given Channel",
			zap.Error(err),
		)
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
		l.Error("error creating Kafka Topic", zap.Error(err))
		return ierrors.Wrap(
			ierrors.New(err).InternalServer(),
			"unable to create kafka topic",
		)
	}
	return nil
}

// Update updates a channel in kafka
func (c *ChannelOperator) Update(
	ctx context.Context,
	context string,
	channel *meta.Channel,
) error {
	logger.Info("trying to update a Channels in Kafka",
		zap.String("channel", channel.Meta.Name),
		zap.String("context", context))
	// updating and creating a new topic is the same thing on kafka
	return c.Create(ctx, context, channel)
}

// Delete deletes a channel from kafka
func (c *ChannelOperator) Delete(
	ctx context.Context,
	context string,
	name string,
) error {
	channel, _ := c.mem.Perm().Channels().Get(context, name)
	topics := []string{toTopic(channel)}
	logger.Info("trying to delete a Channel from Kafka Topics",
		zap.String("channel", name),
		zap.String("context", context))

	_, err := c.k.DeleteTopics(ctx, topics)
	if err != nil {
		logger.Error("error deleting Kafka Topic", zap.Any("error", err))
		return ierrors.Wrap(
			ierrors.New(err).InternalServer(),
			"unable to delete kafka topic",
		)
	}
	return nil
}

type kafkaAdminClient interface {
	DeleteTopics(
		ctx context.Context,
		topics []string,
		options ...kafka.DeleteTopicsAdminOption,
	) (result []kafka.TopicResult, err error)
	CreateTopics(
		ctx context.Context,
		topics []kafka.TopicSpecification,
		options ...kafka.CreateTopicsAdminOption,
	) (result []kafka.TopicResult, err error)
	GetMetadata(
		topic *string,
		allTopics bool,
		timeoutMs int,
	) (*kafka.Metadata, error)
}

type mockAdminClient struct {
}

func (*mockAdminClient) DeleteTopics(
	ctx context.Context,
	topics []string,
	options ...kafka.DeleteTopicsAdminOption,
) (result []kafka.TopicResult, err error) {
	return nil, nil
}

func (*mockAdminClient) CreateTopics(
	ctx context.Context,
	topics []kafka.TopicSpecification,
	options ...kafka.CreateTopicsAdminOption,
) (result []kafka.TopicResult, err error) {
	return nil, nil
}

func (*mockAdminClient) GetMetadata(
	topic *string,
	allTopics bool,
	timeoutMs int,
) (*kafka.Metadata, error) {
	return nil, nil
}
