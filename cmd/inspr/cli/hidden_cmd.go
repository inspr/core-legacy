package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
)

// NewHiddenCmd - hidden subcommand
func NewHiddenCmd() *cobra.Command {
	return cmd.NewCmd("hidden").
		WithExample("should have two arguments", "hidden X Y").
		WithDescription("hidden subcommand").
		Hidden().
		WithCommonFlags().
		WithFlags([]*cmd.Flag{
			{
				Name:          "randomNewTag",
				Usage:         "blablabla",
				Shorthand:     "n",
				Value:         &cmd.InsprOptions.SampleFlagValue,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"hidden"},
			},
		}).
		ExactArgs(2, doHidden)

}

func doHidden(_ context.Context, out io.Writer, strs []string) error {
	fmt.Fprint(out, "hidden hello\n")
	return nil
}
