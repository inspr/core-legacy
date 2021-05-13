package fake

import (
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/utils"
)

// Brokers - mocks the implementation of BrokersMemory interface methods
type Brokers struct {
	fail   error
	broker *brokers.Brokers
}

// GetAll returns an array containing all currently mocked brokers
func (bks *Brokers) GetAll() utils.StringArray {
	return bks.broker.Availible.ToArray()
}

// GetDefault returns the broker mocked as default
func (bks *Brokers) GetDefault() string {
	return bks.broker.Default
}

// Create mocks a new broker on insprd
func (bks *Brokers) Create(broker string, config interface{}) error {
	if bks.fail != nil {
		return bks.fail
	}

	bks.broker.Availible[broker] = true
	return nil
}

// SetDefault sets a previoulsy mocked broker as the fake's default broker
func (bks *Brokers) SetDefault(broker string) error {
	if bks.fail != nil {
		return bks.fail
	}

	bks.broker.Default = broker
	return nil
}
