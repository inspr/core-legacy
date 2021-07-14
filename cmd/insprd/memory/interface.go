package memory

import (
	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
)

// Manager is the interface that allows the management
// of the current state of the cluster. Permiting the
// modification of all Inprd's memory
type Manager interface {
	Tree() tree.Manager
	Brokers() brokers.Manager
}
