package brokers

import "inspr.dev/inspr/pkg/utils"

// Brokers define all Available brokers on Insprd and its default broker
type Brokers struct {
	Default   string
	Available BrokerStatusArray
}

// ChannelBroker associates channel names with their brokers, used to recover data from enviroment
type ChannelBroker struct {
	ChName string
	Broker string
}

// BrokerConfiguration generic interface type
type BrokerConfiguration interface {
	Broker() string
}

// BrokerStatusArray generic status array, used to return brokers data
type BrokerStatusArray map[string]BrokerConfiguration

//Brokers returns an array containing the name of all availible brokers
func (bsa *BrokerStatusArray) Brokers() utils.StringArray {
	arr := utils.StringArray{}
	for k := range *bsa {
		arr = append(arr, k)
	}
	return arr
}
