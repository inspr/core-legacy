package fake

import (
	apimodels "inspr.dev/inspr/pkg/api/models"

	memory "inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/pkg/meta/brokers"
)

// Get returns an array containing all currently mocked brokers
func (bks *BrokersMock) Get() (*apimodels.BrokersDI, error) {
	if bks.fail != nil {
		return nil, bks.fail
	}

	return &apimodels.BrokersDI{
		Available: bks.broker.Available.Brokers(),
		Default:   bks.broker.Default,
	}, nil
}

// GetDefault returns the broker mocked as default
func (bks *BrokersMock) GetDefault() (string, error) {
	if bks.fail != nil {
		return "", bks.fail
	}
	return bks.broker.Default, nil
}

// Create mocks a new broker on insprd
func (bks *BrokersMock) Create(config brokers.BrokerConfiguration) error {
	if bks.fail != nil {
		return bks.fail
	}

	bks.broker.Available[config.Broker()] = config
	return nil
}

// SetDefault sets a previously mocked broker as the fake's default broker
func (bks *BrokersMock) SetDefault(broker string) error {
	if bks.fail != nil {
		return bks.fail
	}

	bks.broker.Default = broker
	return nil
}

// Factory mock of factory interface
func (bks *BrokersMock) Factory() memory.SidecarManager {
	return bks.factory
}

//Configs mock of configuration for broker
func (bks *BrokersMock) Configs(broker string) (brokers.BrokerConfiguration, error) {
	if bks.fail != nil {
		return nil, bks.fail
	}

	return bks.broker.Available[broker], nil
}
