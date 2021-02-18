package main

import (
	"os"

	cli "gitlab.inspr.dev/inspr/core/cmd/inspr/cli"
)

func main() {
	cli.NewInsprCommand(os.Stdout, os.Stderr).Execute()
}
