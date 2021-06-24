package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	dappclient "inspr.dev/inspr/pkg/client"
)

const defaultMOD = 100

func main() {
	var mod int

	// reads from env
	modString, exists := os.LookupEnv("MODULE")
	if !exists {
		mod = defaultMOD
	} else {
		mod, _ = strconv.Atoi(modString)
	}

	// sets up ticker and rand
	rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(200 * time.Millisecond)
	// sets up client for sidecar
	c := dappclient.NewAppClient()
	// channel name
	chName := "primesch1"
	ctx := context.Background()
	fmt.Println("starting...")

	for range ticker.C {
		randNumber := rand.Int() % mod
		fmt.Println("random number -> ", randNumber)
		err := c.WriteMessage(ctx, chName, randNumber)
		fmt.Printf("wrote message to %s\n", chName)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
