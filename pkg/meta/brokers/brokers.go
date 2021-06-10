package brokers

import (
	"github.com/inspr/inspr/pkg/meta/utils"
	pkgutils "github.com/inspr/inspr/pkg/utils"
)

// Brokers define all Available brokers on insprd and its default broker.
type Brokers struct {
	Default   BrokerStatus
	Available utils.StrSet
}

// ChannelBroker associates channels names with their brokers, used to recover data from enviroment
type ChannelBroker struct {
	ChName string
	Broker string
}

// BrokerConfiguration generic interface type
type BrokerConfiguration interface {
	Broker() string
}

// BrokerStatus generiic status type for brokers, used as parameters and returns
type BrokerStatus string

// BrokerStatusArray generic status array, used to return brokers data
type BrokerStatusArray pkgutils.StringArray
