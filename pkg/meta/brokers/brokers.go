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

// BrokerConfiguration generic interface type
type BrokerConfiguration interface{}

// BrokerStatus generiic status type for brokers, used as parameters and returns
type BrokerStatus string

// BrokerStatusArray generic status array, used to return brokers data
type BrokerStatusArray pkgutils.StringArray
