// THIS IS THE MASTER
package main

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/api"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
	kafka "gitlab.inspr.dev/inspr/core/cmd/insprd/operators/kafka"
)

func main() {
	memoryManager := tree.GetTreeMemory()
	channelOperator, err := kafka.NewKafkaOperator()
	if err != nil {
		panic(err)
	}

	api.Run(memoryManager, channelOperator)
}
