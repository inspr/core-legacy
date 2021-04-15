package main

import (
	"os"

	"github.com/inspr/inspr/cmd/insprd/api"
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/memory/tree"
	"github.com/inspr/inspr/cmd/insprd/operators"
	kafka "github.com/inspr/inspr/cmd/insprd/operators/kafka"
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
