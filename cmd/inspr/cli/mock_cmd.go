package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// NewMockCmd - mock subcommand
func NewMockCmd() *cobra.Command {
	return NewCmd("mock").
		WithDescription("mocks the usage of a subcommand").
		WithCommonFlags().
		NoArgs(doMock)
}

func doMock(_ context.Context, out io.Writer) error {
	fmt.Fprint(out, "mock hello\n")
	return nil
}
