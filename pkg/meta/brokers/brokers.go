package brokers

import "github.com/inspr/inspr/pkg/meta/utils"

// Brokers define all availible brokers on insprd and its default broker.
type Brokers struct {
	Default   string
	Availible utils.StrSet
}
