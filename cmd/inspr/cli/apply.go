package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
)

// NewApplyCmd - mock subcommand
func NewApplyCmd() *cobra.Command {
	return cmd.NewCmd("apply").
		WithDescription("mocks the usage of a subcommand").
		WithCommonFlags().
		NoArgs(doMock)
}

func doApply(_ context.Context, out io.Writer) error {
	fmt.Fprint(out, "mock hello\n")
	return nil
}
