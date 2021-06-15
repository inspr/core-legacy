package brokers

import (
	"github.com/inspr/inspr/cmd/sidecars"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/sidecars/models"
	"github.com/inspr/inspr/pkg/utils"
	"go.uber.org/zap"
)

// GetAll returns an array containing all currently configured brokers
func (bmm *BrokerMemoryManager) GetAll() (utils.StringArray, error) {
	mem, err := bmm.get()
	if err != nil {
		return nil, err
	}
	return mem.Available.Brokers(), nil
}

// GetDefault returns the broker configured as default
func (bmm *BrokerMemoryManager) GetDefault() (string, error) {
	mem, err := bmm.get()
	if err != nil {
		return "", err
	}
	return mem.Default, nil
}

func (bmm *BrokerMemoryManager) get() (*brokers.Brokers, error) {
	if bmm.broker == nil {
		return nil, ierrors.NewError().Message("broker status memory is empty").Build()
	}
	return bmm.broker, nil
}

// Create configures a new broker on insprd
func (bmm *BrokerMemoryManager) Create(config brokers.BrokerConfiguration) error {
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
func (bmm *BrokerMemoryManager) SetDefault(broker string) error {
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
func (bmm *BrokerMemoryManager) Factory() SidecarManager {
	return bmm.factory
}

//Configs returns the configurations for a given broker
func (bmm *BrokerMemoryManager) Configs(broker string) (brokers.BrokerConfiguration, error) {
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
