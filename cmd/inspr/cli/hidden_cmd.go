package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
)

var randomNewTag string

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
				Value:         &randomNewTag,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"hidden"},
			},
		}).
		ExactArgs(2, doHidden)

}

func doHidden(_ context.Context, out io.Writer, strs []string) error {
	fmt.Fprintf(out, "hidden hello -> %v\n", randomNewTag)
	return nil
}
