package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// NewConfigChangeCmd - responsible for changing the values of the inspr's viper config
func NewConfigChangeCmd() *cobra.Command {
	return cmd.NewCmd("config <key> <value>").
		WithDescription("Change the values stored in the inspr config").
		WithExample("Changing IP config", "config serverip http://127.0.0.1:8080").
		WithExample("Changing scope config", "config scope app1.app2").
		AddSubCommand(NewListConfig()).
		ExactArgs(2, doConfigChange)
}

// NewListConfig - config subcommand that shows all existant variables in the config
func NewListConfig() *cobra.Command {
	return cmd.NewCmd("list").
		WithDescription("See the list of configuration variables and their current values").
		WithExample("type", "config list").
		NoArgs(doListConfig)
}

func doConfigChange(_ context.Context, out io.Writer, args []string) error {
	key := args[0]
	value := args[1]

	// key doesn't exist
	if !cliutils.ExistsKey(key) {
		errMsg := "error: key inserted does not exist in the inspr config"
		fmt.Fprintln(out, errMsg)
		printExistingKeys(out)
		return ierrors.NewError().Message(errMsg).Build()
	}

	// updates
	if err := cliutils.ChangeViperValues(key, value); err != nil {
		return err
	}

	fmt.Fprintf(out, "Success: inspr config [%v] changed to '%v'\n", key, value)
	return nil
}

func doListConfig(_ context.Context, out io.Writer) error {
	printExistingKeys(out)
	return nil
}

func printExistingKeys(out io.Writer) {
	fmt.Fprintln(out, "Available configurations: ")
	for _, key := range cliutils.ExistingKeys() {
		value := viper.GetString(key)
		value = "\"" + value + "\""
		fmt.Fprintf(out, "- %v: %v\n", key, value)
	}
}
