package brokers

import (
	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/sidecars"
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta/brokers"
	"inspr.dev/inspr/pkg/sidecars/models"
)

// Get returns the brokers configured data
func (bmm *brokerMemoryManager) Get() (*apimodels.BrokersDI, error) {
	mem, err := bmm.get()
	if err != nil {
		return nil, err
	}
	return &apimodels.BrokersDI{
		Available: mem.Available.Brokers(),
		Default:   mem.Default,
	}, nil
}

func (bmm *brokerMemoryManager) get() (*brokers.Brokers, error) {
	if bmm.broker == nil {
		return nil, ierrors.NewError().Message("broker status memory is empty").Build()
	}
	return bmm.broker, nil
}

// Create configures a new broker on insprd
func (bmm *brokerMemoryManager) Create(config brokers.BrokerConfiguration) error {
	logger.Info("creating new broker")
	mem, err := bmm.get()
	if err != nil {
		return err
	}

	broker := config.Broker()
	logger.Debug("broker to be created",
		zap.String("broker", broker),
		zap.Any("configs", config))

	if _, ok := mem.Available[broker]; ok {
		return ierrors.NewError().Message("broker %s is already configured on memory", broker).Build()
	}

	var factory models.SidecarFactory
	switch broker {
	case brokers.Kafka:
		obj, _ := config.(*sidecars.KafkaConfig)
		factory = sidecars.KafkaToDeployment(*obj)
	default:
		return ierrors.NewError().Message("broker %s is not supported", broker).Build()
	}

	logger.Debug("subscribing broker to sidecar factory")
	err = bmm.Factory().Subscribe(broker, factory)
	if err != nil {
		logger.Error("unable to subscribe broker")
		return err
	}

	mem.Available[broker] = config
	if mem.Default == "" {
		bmm.SetDefault(broker)
	}
	return nil
}

// SetDefault sets a previously configured broker as insprd's default broker
func (bmm *brokerMemoryManager) SetDefault(broker string) error {
	logger.Debug("setting new default broker",
		zap.String("broker", broker))
	mem, err := bmm.get()
	if err != nil {
		return err
	}

	if _, ok := mem.Available[broker]; !ok {
		return ierrors.NewError().Message("broker %s is not configured on memory", broker).Build()
	}

	mem.Default = broker
	return nil
}

// Factory provides the struct implementation for Sidecarfactory
func (bmm *brokerMemoryManager) Factory() SidecarManager {
	return bmm.factory
}

//Configs returns the configurations for a given broker
func (bmm *brokerMemoryManager) Configs(broker string) (brokers.BrokerConfiguration, error) {
	logger.Info("getting config for broker sidecar",
		zap.String("broker", broker))

	mem, err := bmm.get()
	if err != nil {
		return nil, err
	}

	config, ok := mem.Available[broker]
	if !ok {
		return nil, ierrors.NewError().Message("broker %s is not configured on memory", broker).Build()
	}

	return config, nil
}
