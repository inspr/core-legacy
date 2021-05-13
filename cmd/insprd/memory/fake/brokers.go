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

func (bks *Brokers) GetAll() utils.StringArray {
	return bks.broker.Availible.ToArray()
}

func (bks *Brokers) GetDefault() string {
	return bks.broker.Default
}

func (bks *Brokers) Create(broker string, config interface{}) error {
	if bks.fail != nil {
		return bks.fail
	}

	bks.broker.Availible[broker] = true
	return nil
}

func (bks *Brokers) SetDefault(broker string) error {
	if bks.fail != nil {
		return bks.fail
	}

	bks.broker.Default = broker
	return nil
}
