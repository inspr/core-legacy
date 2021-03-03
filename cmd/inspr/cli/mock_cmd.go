package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"

	"gitlab.inspr.dev/inspr/core/pkg/cmd"
)

// NewMockCmd - mock subcommand
func NewMockCmd() *cobra.Command {
	return cmd.NewCmd("mock").
		WithDescription("mocks the usage of a subcommand").
		WithCommonFlags().
		AddSubCommand(NewHelloCmd()).
		NoArgs(doMock)
}

func doMock(_ context.Context) error {
	out := cliutils.GetCliOut()
	value := cliutils.GetConfiguredScope()
	fmt.Fprintln(out, value)

	port := cliutils.GetConfiguredServerIP()
	fmt.Fprintln(out, port)

	return nil
}
