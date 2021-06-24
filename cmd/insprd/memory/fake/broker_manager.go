package fake

import (
	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	metabroker "inspr.dev/inspr/pkg/meta/brokers"
)

// BrokersMock is the struct with the necessary implementations
// to mock the interface used to manage brokers
type BrokersMock struct {
	fail    error
	broker  *metabroker.Brokers
	factory *Factory
}

// MockBrokerMemory mock exported with propagated error through the functions
func MockBrokerMemory(failErr error) brokers.Manager {
	return &BrokersMock{
		fail: failErr,
		broker: &metabroker.Brokers{
			Default: "default_mock",
			Available: metabroker.BrokerStatusArray{
				"default_mock": nil,
			},
		},
	}
}
