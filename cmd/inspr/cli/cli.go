package cli

import (
	"io"
	"os"

	"github.com/inspr/inspr/pkg/cmd"
	cliutils "github.com/inspr/inspr/pkg/cmd/utils"

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
		).
		Version(version).
		WithLongDescription("main command of the inspr cli, to see the full list of subcommands existent please use 'inspr help'").
		Super()

	rootCmd.PersistentPreRunE = mainCmdPreRun

	// root persistentFlags
	return rootCmd
}

func mainCmdPreRun(cm *cobra.Command, args []string) error {
	cm.Root().SilenceUsage = true
	// viper defaults values or reads from the config location
	cliutils.InitViperConfig()

	homeDir, _ := os.UserHomeDir()

	if err := cliutils.ReadViperConfig(homeDir); err != nil {
		return err
	}

	return nil
}
