package fake

import (
	"github.com/inspr/inspr/cmd/insprd/memory/brokers"
	metabroker "github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/meta/utils"
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
	availibe, _ := utils.MakeStrSet([]string{"default_mock"})
	return &BrokersMock{
		fail: failErr,
		broker: &metabroker.Brokers{
			Default:   "default_mock",
			Available: availibe,
		},
	}
}
