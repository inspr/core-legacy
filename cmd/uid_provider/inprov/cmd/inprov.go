package cmd

import (
	"gitlab.inspr.dev/inspr/core/cmd/uid_provider/inprov/client"
	build "gitlab.inspr.dev/inspr/core/pkg/cmd"
)

var MainCommand = build.NewCmd("inprov <subcommand>").AddSubCommand(
	createUserCmd,
	deleteUserCmd,
	loginCmd,
).Super()

func init() {
	cl = client.NewClient()
}
