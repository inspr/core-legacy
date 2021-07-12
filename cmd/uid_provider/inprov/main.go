package main

import "inspr.dev/inspr/cmd/uid_provider/inprov/cmd"

func main() {
	cmd.MainCommand.Root().SilenceUsage = true
	cmd.MainCommand.Execute()
}
