package memory

import (
	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
)

type Manager interface {
	Tree() tree.Manager
	Brokers() brokers.Manager
}
