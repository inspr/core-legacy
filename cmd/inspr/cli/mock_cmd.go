package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func doMock(_ context.Context, out io.Writer) error {
	value := viper.Get("mock")
	fmt.Fprintln(out, value)

	port := viper.Get("port")
	fmt.Fprintln(out, port)

	return nil
}
