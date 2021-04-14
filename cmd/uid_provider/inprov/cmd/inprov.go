package cmd

import (
	"github.com/inspr/inspr/cmd/uid_provider/inprov/client"
	build "github.com/inspr/inspr/pkg/cmd"
)

var MainCommand = build.NewCmd("inprov <subcommand>").AddSubCommand(
	createUserCmd,
	deleteUserCmd,
	loginCmd,
).Super()

func init() {
	cl = client.NewClient()
}
