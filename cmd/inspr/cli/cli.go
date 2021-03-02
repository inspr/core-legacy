package cli

import (
	"io"

	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"

	"github.com/spf13/cobra"
)

// NewInsprCommand - returns a root command associated with inspr cli
func NewInsprCommand(out, err io.Writer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "inspr",
		Short:         "main command of the inspr cli",
		Long:          `main command of the inspr cli, to see the full list of subcommands existant please use 'inspr help'`,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.Root().SilenceUsage = true

			// viper defaults values or reads from the config location
			cliutils.InitViperConfig()

			if err := cliutils.ReadViperConfig(); err != nil {
				return err
			}

			return nil
		},
	}

	// other commmands
	rootCmd.AddCommand(NewMockCmd())
	rootCmd.AddCommand(NewHiddenCmd())
	rootCmd.AddCommand(NewGetCmd())
	rootCmd.AddCommand(NewDeleteCmd())

	rootCmd.AddCommand(NewApplyCmd())
	rootCmd.AddCommand(NewDescribeCmd())

	rootCmd.AddCommand(NewConfigChangeCmd())
	// root persistentFlags
	return rootCmd
}
