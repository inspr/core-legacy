package cli

import (
	"io"

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
			return nil
		},
	}

	// other commmands
	rootCmd.AddCommand(NewMockCmd())
	rootCmd.AddCommand(NewHiddenCmd())
	rootCmd.AddCommand(NewApplyCmd())
	rootCmd.AddCommand(NewDescribeCmd())
	// root persistentFlags
	return rootCmd
}
