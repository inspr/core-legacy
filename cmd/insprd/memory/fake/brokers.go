package fake

import (
	memory "github.com/inspr/inspr/cmd/insprd/memory/brokers"
	"github.com/inspr/inspr/pkg/meta/brokers"
)

// GetAll returns an array containing all currently mocked brokers
func (bks *BrokersMock) GetAll() brokers.BrokerStatusArray {
	return brokers.BrokerStatusArray(bks.broker.Available.ToArray())
}

// GetDefault returns the broker mocked as default
func (bks *BrokersMock) GetDefault() brokers.BrokerStatus {
	return brokers.BrokerStatus(bks.broker.Default)
}

// Create mocks a new broker on insprd
func (bks *BrokersMock) Create(broker brokers.BrokerStatus, config brokers.BrokerConfiguration) error {
	if bks.fail != nil {
		return bks.fail
	}

	bks.broker.Available[string(broker)] = true
	return nil
}

// SetDefault sets a previoulsy mocked broker as the fake's default broker
func (bks *BrokersMock) SetDefault(broker brokers.BrokerStatus) error {
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
