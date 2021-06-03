package main

import (
	"os"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/memory/brokers"
	"github.com/inspr/inspr/cmd/insprd/memory/tree"
	"github.com/inspr/inspr/cmd/insprd/operators"
	"github.com/inspr/inspr/cmd/sidecars"
	"github.com/inspr/inspr/pkg/api"
	"github.com/inspr/inspr/pkg/auth"
	jwtauth "github.com/inspr/inspr/pkg/auth/jwt"
	authmock "github.com/inspr/inspr/pkg/auth/mocks"
	metabrokers "github.com/inspr/inspr/pkg/meta/brokers"
)

func main() {
	var memoryManager memory.Manager
	var brokerManager brokers.Manager
	var operator operators.OperatorInterface
	var authenticator auth.Auth
	var err error

	if _, ok := os.LookupEnv("DEBUG"); ok {
		authenticator = authmock.NewMockAuth(nil)
		brokerManager = brokers.GetBrokerMemory()
		memoryManager = tree.GetTreeMemory()
		operator, err = operators.NewOperator(memoryManager, authenticator, brokerManager)
		if err != nil {
			panic(err)
		}
	} else {
		pubKey, err := auth.GetPublicKey()
		if err != nil {
			panic(err)
		}
		authenticator = jwtauth.NewJWTauth(pubKey)
		brokerManager = brokers.GetBrokerMemory()
		memoryManager = tree.GetTreeMemory()
		operator, err = operators.NewOperator(memoryManager, authenticator, brokerManager)
		if err != nil {
			panic(err)
		}
	}
	config := sidecars.KafkaConfig{
		BootstrapServers: "kafka.default.svc:9092",
		AutoOffsetReset:  "earliest",
		KafkaInsprAddr:   "http://localhost",
		SidecarImage:     "gcr.io/red-inspr/inspr/sidecar/kafka:latest",
	}

	brokerManager.Create(metabrokers.BrokerStatus(metabrokers.Kafka), config)

	api.Run(memoryManager, operator, authenticator, brokerManager)
}
