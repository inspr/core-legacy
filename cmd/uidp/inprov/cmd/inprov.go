package cmd

import (
	"inspr.dev/inspr/cmd/uidp/inprov/client"
	build "inspr.dev/inspr/pkg/cmd"
)

// MainCommand is the main command for the inspr uid provider CLI, aka inprov.
var MainCommand = build.NewCmd("inprov <subcommand>").AddSubCommand(
	createUserCmd,
	deleteUserCmd,
	loginCmd,
).Super()

func init() {
	cl = client.NewClient()
}
