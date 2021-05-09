package cli

import (
	"fmt"
	"io"

	"github.com/inspr/inspr/pkg/cmd"
	"github.com/inspr/inspr/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

// NewInsprCommand - returns a root command associated with inspr cli
func NewInsprCommand(out, err io.Writer, version string) *cobra.Command {
	rootCmd := cmd.NewCmd("inspr").
		WithDescription("main command of the inspr cli").
		WithCommonFlags().
		AddSubCommand(NewGetCmd(),
			NewDeleteCmd(),
			NewApplyCmd(),
			NewDescribeCmd(),
			NewConfigChangeCmd(),
			authCommand,
			initCommand,
		).
		Version(version).
		WithLongDescription("main command of the inspr cli, to see the full list of subcommands existent please use 'inspr help'").
		Super()

	rootCmd.PersistentPreRunE = mainCmdPreRun

	// root persistentFlags
	return rootCmd
}

func mainCmdPreRun(cm *cobra.Command, args []string) error {
	if cm.Name() == "init" {
		return nil
	}
	cm.Root().SilenceErrors = true
	cm.Root().SilenceUsage = true
	utils.InitViperConfig()
	// viper defaults values or reads from the config location
	var err error
	if cmd.InsprOptions.Config == "" {
		err = utils.ReadDefaultConfig()
	} else {
		err = utils.ReadConfigFromFile(cmd.InsprOptions.Config)
	}
	if err != nil {
		fmt.Fprintf(utils.GetCliOutput(), "Invalid config file! Did you run inspr init?")
	}
	return nil
}
