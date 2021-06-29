package brokers

import (
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/meta/brokers"
)

type BrokerMockManager struct {
	broker *brokers.Brokers
	err    error
}

func (mock *BrokerMockManager) Get() (*apimodels.BrokersDI, error) {
	if mock.err != nil {
		return nil, mock.err
	}
	return &apimodels.BrokersDI{
		Available: mock.broker.Available.Brokers(),
		Default:   mock.broker.Default,
	}, nil
}

// Create configures a new broker on insprd
func (mock *BrokerMockManager) Create(config brokers.BrokerConfiguration) error {
	if mock.err != nil {
		return mock.err
	}
	mock.broker.Available[config.Broker()] = config
	if mock.broker.Default == "" {
		mock.SetDefault(config.Broker())
	}
	return nil
}

// SetDefault sets a previously configured broker as insprd's default broker
func (mock *BrokerMockManager) SetDefault(broker string) error {
	if mock.err != nil {
		return mock.err
	}
	mock.broker.Default = broker
	return nil
}

// Factory provides the struct implementation for Sidecarfactory
func (mock *BrokerMockManager) Factory() SidecarManager {
	return nil
}

//Configs returns the configurations for a given broker
func (mock *BrokerMockManager) Configs(broker string) (brokers.BrokerConfiguration, error) {
	if mock.err != nil {
		return nil, mock.err
	}
	config, ok := mock.broker.Available[broker]
	if !ok {
		return nil, mock.err
	}
	return config, nil
}
