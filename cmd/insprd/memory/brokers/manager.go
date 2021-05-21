package brokers

import (
	"github.com/inspr/inspr/pkg/meta/brokers"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
)

// BrokerManager implements broker's Manager interface,
// allows for management of the system's message brokers

// BrokerMemoryManager implements the methods described by the BrokersInterface
type BrokerMemoryManager struct {
	factory SidecarManager
	broker  *brokers.Brokers
}

var brokerMemory Manager

// GetBrokerMemory allows for connection with BrokersManager sigleton
func GetBrokerMemory() Manager {
	if brokerMemory == nil {
		brokerMemory = &BrokerMemoryManager{
			broker: &brokers.Brokers{
				Available: make(metautils.StrSet),
			},
			factory: &AbstractBrokerFactory{},
		}
	}
	return brokerMemory
}

// ResetBrokerMemory makes the BrokersManagers singleton points to nil,
// used only in tests
func ResetBrokerMemory() {
	brokerMemory = nil
}