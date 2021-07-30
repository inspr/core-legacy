package main

import "inspr.dev/inspr/cmd/uidp/inprov/cmd"

func main() {
	cmd.MainCommand.Root().SilenceUsage = true
	cmd.MainCommand.Execute()
}
