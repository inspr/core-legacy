package brokers

type BrokersManager struct {
}

var brokerMemory Manager

// GetBrokerMemory
func GetBrokerMemory() Manager {
	if brokerMemory == nil {
		brokerMemory = &BrokersManager{}
	}
	return brokerMemory
}

// Brokers
func (bm *BrokersManager) Brokers() BrokerInterface {
	return &BrokersMemoryManager{}
}
