package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
)

// NewHelloCmd - hello subcommand of a subcommand
func NewHelloCmd() *cobra.Command {
	return cmd.NewCmd("hello").
		WithDescription("HELLO WORLD").
		NoArgs(doHello)
}

func doHello(_ context.Context) error {
	out := cliutils.GetCliOut()
	fmt.Fprint(out, "sub of a subcommand hello\n")
	return nil
}
