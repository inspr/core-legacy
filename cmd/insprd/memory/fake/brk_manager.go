package fake

import (
	"github.com/inspr/inspr/cmd/insprd/memory/brokers"
	metabroker "github.com/inspr/inspr/pkg/meta/brokers"
)

// BrkManager is the api struct with the necessary implementations
// to mock the interface used to manage brokers
type BrkManager struct {
	brokers Brokers
}

// MockBrokerManager mock exported with propagated error through the functions
func MockBrokerManager(failErr error) brokers.Manager {
	return &BrkManager{
		brokers: Brokers{
			fail:   failErr,
			broker: &metabroker.Brokers{},
		},
	}
}

// Brokers mock of broker interface
func (bm *BrkManager) Brokers() brokers.BrokerInterface {
	return &bm.brokers
}
