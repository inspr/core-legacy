// THIS IS THE MASTER
package main

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/api"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
	kafka "gitlab.inspr.dev/inspr/core/cmd/insprd/operators/kafka/channels"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators/kafka/nodes"
)

func main() {
	memoryManager := tree.GetTreeMemory()
	channelOperator, err := kafka.NewOperator()
	if err != nil {
		panic(err)
	}
	nodeOperator := nodes.NewOperator()
	api.Run(memoryManager, nodeOperator, channelOperator)
}
