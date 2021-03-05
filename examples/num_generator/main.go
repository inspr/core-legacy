package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	kafkasc "gitlab.inspr.dev/inspr/core/cmd/sidecars/kafka/client"
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
	writer, err := kafkasc.NewWriter(false)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ticker.C:
			writer.WriteMessage("ch1", (rand.Int() % mod))
		}
	}
}
