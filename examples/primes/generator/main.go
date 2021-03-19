package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	dappclient "gitlab.inspr.dev/inspr/core/pkg/client"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
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
	ticker := time.NewTicker(2 * time.Second)
	rand.Seed(time.Now().UnixNano())

	// sets up client for sidecar
	c := dappclient.NewAppClient()
	// channelName
	chName := "primes_ch1"
	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			randNumber := rand.Int() % mod
			fmt.Println("random number -> ", randNumber)
			newMsg := models.Message{
				Data: randNumber,
			}

			err := c.WriteMessage(ctx, chName, newMsg)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
