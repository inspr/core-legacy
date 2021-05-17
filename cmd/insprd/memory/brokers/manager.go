package brokers

// BrokerManager implements broker's Manager interface,
// allows for management of the system's message brokers
type BrokerManager struct {
}

var brokerMemory Manager

// GetBrokerMemory allows for connection with BrokersManager sigleton
func GetBrokerMemory() Manager {
	if brokerMemory == nil {
		brokerMemory = &BrokerManager{}
	}
	return brokerMemory
}

// Brokers provides access to brokers memory
func (bm *BrokerManager) Brokers() BrokerInterface {
	return &BrokerMemoryManager{
		factory: &AbstractBrokerFactory{},
	}
}
