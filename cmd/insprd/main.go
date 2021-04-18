package main

import (
	"os"

	"inspr.dev/inspr/cmd/insprd/api"
	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
	"inspr.dev/inspr/cmd/insprd/operators"
	kafka "inspr.dev/inspr/cmd/insprd/operators/kafka"
)

func main() {
	var memoryManager memory.Manager
	var operator operators.OperatorInterface

	var err error
	if _, ok := os.LookupEnv("DEBUG"); ok {
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
