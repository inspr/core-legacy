package brokers

import (
	"sync"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/meta/brokers"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "broker-memory")))
	// logger = zap.NewNop()
}

// BrokerManager implements broker's Manager interface,
// allows for management of the system's message brokers

// brokerMemoryManager implements the methods described by the BrokersInterface
type brokerMemoryManager struct {
	factory   SidecarManager
	broker    *brokers.Brokers
	availible sync.Mutex
	def       sync.Mutex
}

var brokerMemory *brokerMemoryManager

// GetBrokerMemory allows for connection with BrokersManager sigleton
func GetBrokerMemory() Manager {
	if brokerMemory == nil {
		brokerMemory = &brokerMemoryManager{
			broker: &brokers.Brokers{
				Available: make(brokers.BrokerStatusArray),
			},
			factory: &AbstractBrokerFactory{},
		}
	}
	return brokerMemory
}
