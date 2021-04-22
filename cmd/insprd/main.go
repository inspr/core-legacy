package main

import (
	"os"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/memory/tree"
	"github.com/inspr/inspr/cmd/insprd/operators"
	kafka "github.com/inspr/inspr/cmd/insprd/operators/kafka"
	"github.com/inspr/inspr/pkg/api"
	"github.com/inspr/inspr/pkg/auth"
	jwtauth "github.com/inspr/inspr/pkg/auth/jwt"
)

func main() {
	var memoryManager memory.Manager
	var operator operators.OperatorInterface
	var authenticator auth.Auth
	var err error

	pubKey, err := auth.GetPublicKey()
	if err != nil {
		panic(err)
	}

	if _, ok := os.LookupEnv("DEBUG"); ok {
		authenticator = jwtauth.NewJWTauth(pubKey)
		memoryManager = tree.GetTreeMemory()
		operator, err = kafka.NewKafkaOperator(memoryManager)
		if err != nil {
			panic(err)
		}
	} else {
		authenticator = jwtauth.NewJWTauth(pubKey)
		memoryManager = tree.GetTreeMemory()
		operator, err = kafka.NewKafkaOperator(memoryManager)
		if err != nil {
			panic(err)
		}
	}

	api.Run(memoryManager, operator, authenticator)
}
