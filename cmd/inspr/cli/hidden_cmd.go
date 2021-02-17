package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// NewHiddenCmd - hidden subcommand
func NewHiddenCmd() *cobra.Command {
	return NewCmd("hidden").
		WithDescription("hidden subcommand").
		Hidden().
		WithCommonFlags().
		NoArgs(doHidden)
}

func doHidden(_ context.Context, out io.Writer) error {
	fmt.Fprint(out, "hidden hello\n")
	return nil
}
