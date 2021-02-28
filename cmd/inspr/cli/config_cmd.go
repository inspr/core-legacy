package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// NewConfigChangeCmd - responsible for changing the values of the inspr's viper config
func NewConfigChangeCmd() *cobra.Command {
	return cmd.NewCmd("config").
		WithDescription("Used to changed the values stored in the inspr config").
		WithExample("how to use", "config <key> <value>").
		ExactArgs(2, doConfigChange)
}

func doConfigChange(_ context.Context, out io.Writer, args []string) error {
	key := args[0]
	value := args[1]

	// key doesn't exist
	if !existsKey(key) {
		errMsg := "key inserted does not exist in the inspr config"
		fmt.Fprintln(out, errMsg)
		fmt.Fprintln(out, "existing keys")
		fmt.Fprintln(out, existingKeys())
		return ierrors.NewError().Message(errMsg).Build()
	}

	// updates
	if err := changeViperValues(key, value); err != nil {
		return err
	}

	return nil
}