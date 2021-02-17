// THIS IS THE MASTER
package main

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/api"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
)

func main() {
	memoryManager := tree.GetTreeMemory()
	api.Run(memoryManager)
}
