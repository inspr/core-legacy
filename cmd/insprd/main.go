package main

import (
	"os"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/fake"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	kafka "gitlab.inspr.dev/inspr/core/cmd/insprd/operators/kafka"
)

func main() {
	var memoryManager memory.Manager
	var operator operators.OperatorInterface

	var err error
	if _, ok := os.LookupEnv("DEBUG"); ok {
		memoryManager = fake.MockMemoryManager(nil)
		memoryManager = tree.GetTreeMemory()
		operator, err = kafka.NewKafkaOperator(memoryManager)
		if err != nil {
			panic(err)
		}
	} else {
		memoryManager = tree.GetTreeMemory()
		operator, err = kafka.NewKafkaOperator(memoryManager)
		if err != nil {
			panic(err)
		}
	}

	api.Run(memoryManager, operator)
}
