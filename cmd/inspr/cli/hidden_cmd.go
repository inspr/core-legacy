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
		WithFlags([]*Flag{
			{
				Name:          "randomNewTag",
				Usage:         "blablabla",
				Shorthand:     "n",
				Value:         &InsprOptions.sampleFlagValue,
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
