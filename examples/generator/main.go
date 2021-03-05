package main

import (
	"context"
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
	chName := "ch1"

	for {
		select {
		case <-ticker.C:
			c.WriteMessage(context.Background(), chName, models.Message{
				Data: (rand.Int() % mod),
			})
		}
	}
}
