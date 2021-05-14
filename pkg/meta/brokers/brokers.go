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

type BrokerConfiguration interface{}
type BrokerStatus string
type BrokerStatusArray pkgutils.StringArray
