package fake

import (
	"github.com/inspr/inspr/cmd/insprd/memory/brokers"
	metabroker "github.com/inspr/inspr/pkg/meta/brokers"
)

type BrkManager struct {
	brokers Brokers
}

func MockBrokerManager(failErr error) brokers.Manager {
	return &BrkManager{
		brokers: Brokers{
			fail:   failErr,
			broker: &metabroker.Brokers{},
		},
	}
}

func (bm *BrkManager) Brokers() brokers.BrokerInterface {
	return &bm.brokers
}
