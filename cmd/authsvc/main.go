package main

import (
	"fmt"
	"os"

	"gitlab.inspr.dev/inspr/core/cmd/authsvc/api"
)

func main() {
	key := os.Getenv("JWT_PRIVATE_KEY")
	fmt.Println(key)
	api.Run()
}
