package main

import (
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/memory/tree"
	"github.com/inspr/inspr/pkg/api"
	"github.com/inspr/inspr/pkg/auth"
	jwtauth "github.com/inspr/inspr/pkg/auth/jwt"
)

func main() {
	var memoryManager memory.Manager
	var authenticator auth.Auth
	var err error

	pubKey, err := auth.GetPublicKey()
	if err != nil {
		panic(err)
	}

	authenticator = jwtauth.NewJWTauth(pubKey)
	memoryManager = tree.GetTreeMemory()

	api.Run(memoryManager, authenticator)
}
