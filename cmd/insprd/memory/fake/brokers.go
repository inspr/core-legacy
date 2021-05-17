package fake

import (
	memory "github.com/inspr/inspr/cmd/insprd/memory/brokers"
	"github.com/inspr/inspr/pkg/meta/brokers"
)

// Brokers - mocks the implementation of BrokersMemory interface methods
type Brokers struct {
	fail    error
	broker  *brokers.Brokers
	factory *Factory
}

// GetAll returns an array containing all currently mocked brokers
func (bks *Brokers) GetAll() brokers.BrokerStatusArray {
	return brokers.BrokerStatusArray(bks.broker.Available.ToArray())
}

// GetDefault returns the broker mocked as default
func (bks *Brokers) GetDefault() brokers.BrokerStatus {
	return brokers.BrokerStatus(bks.broker.Default)
}

// Create mocks a new broker on insprd
func (bks *Brokers) Create(broker brokers.BrokerStatus, config brokers.BrokerConfiguration) error {
	if bks.fail != nil {
		return bks.fail
	}

	bks.broker.Available[string(broker)] = true
	return nil
}

// SetDefault sets a previoulsy mocked broker as the fake's default broker
func (bks *Brokers) SetDefault(broker brokers.BrokerStatus) error {
	if bks.fail != nil {
		return bks.fail
	}

	bks.broker.Default = broker
	return nil
}

// Factory mock of factory interface
func (bks *Brokers) Factory() memory.SidecarInterface {
	return bks.factory
}
