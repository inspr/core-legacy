package fake

import (
	"github.com/inspr/inspr/cmd/insprd/memory/brokers"
	metabroker "github.com/inspr/inspr/pkg/meta/brokers"
)

// BrokersMock is the struct with the necessary implementations
// to mock the interface used to manage brokers
type BrokersMock struct {
	fail    error
	broker  *metabroker.Brokers
	factory *Factory
}

// MockBrokerManager mock exported with propagated error through the functions
func MockBrokerManager(failErr error) brokers.Manager {
	return &BrokersMock{
		fail:   failErr,
		broker: &metabroker.Brokers{},
	}
}
