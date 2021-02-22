package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
)

// NewHelloCmd - hello subcommand of a subcommand
func NewHelloCmd() *cobra.Command {
	return cmd.NewCmd("hello").
		WithDescription("HELLO WORLD").
		WithCommonFlags().
		NoArgs(doHello)
}

func doHello(_ context.Context, out io.Writer) error {
	fmt.Fprint(out, "sub of a subcommand hello\n")
	return nil
}
