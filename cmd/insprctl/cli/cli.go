package cli

import (
	"fmt"
	"io"

	"inspr.dev/inspr/pkg/cmd"
	"inspr.dev/inspr/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

// NewInsprCommand - returns a root command associated with inspr cli
func NewInsprCommand(out, err io.Writer, version string) *cobra.Command {
	rootCmd := cmd.NewCmd("insprctl").
		WithDescription("main command of the insprctl cli").
		WithCommonFlags().
		AddSubCommand(NewGetCmd(),
			NewDeleteCmd(),
			NewApplyCmd(),
			NewDescribeCmd(),
			NewConfigChangeCmd(),
			completionCmd,
			NewClusterCommand(),
			NewBrokerCmd(),
			initCommand,
		).
		Version(version).
		WithLongDescription("main command of the inspr cli, to see the full list of subcommands existent please use 'insprctl help'").
		Super()

	rootCmd.PersistentPreRunE = mainCmdPreRun

	// root persistentFlags
	return rootCmd
}

func mainCmdPreRun(cm *cobra.Command, args []string) error {
	if cm.Name() == "completion" {
		return nil
	}
	if cm.Name() == "init" && cm.Parent().Name() == "insprctl" {
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
		fmt.Fprintln(
			utils.GetCliOutput(),
			"Invalid config file! Did you run insprctl init?",
		)
	}
	return err
}
