package fake

import (
	memory "github.com/inspr/inspr/cmd/insprd/memory/brokers"
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/utils"
)

// GetAll returns an array containing all currently mocked brokers
func (bks *BrokersMock) GetAll() (utils.StringArray, error) {
	if bks.fail != nil {
		return nil, bks.fail
	}

	return bks.broker.Available.Brokers(), nil
}

// GetDefault returns the broker mocked as default
func (bks *BrokersMock) GetDefault() (string, error) {
	if bks.fail != nil {
		return "", bks.fail
	}
	return bks.broker.Default, nil
}

// Create mocks a new broker on insprd
func (bks *BrokersMock) Create(broker string, config brokers.BrokerConfiguration) error {
	if bks.fail != nil {
		return bks.fail
	}

	bks.broker.Available[string(broker)] = config
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
