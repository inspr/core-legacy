package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			initViper()

			if err := viper.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					if createErr := createConfig(); createErr != nil {
						err = createErr
						fmt.Fprintln(out, err.Error())
					} else {
						err = nil
					}
				} else {
					fmt.Fprintln(out, err.Error())
				}
				return err
			}

			err := viper.Unmarshal(&conf)
			if err != nil {
				fmt.Fprintln(out, "Error unmarshaling the config file")
				fmt.Fprintln(out, err.Error())
			}
			return nil
		},
	}

	// other commmands
	rootCmd.AddCommand(NewMockCmd())
	rootCmd.AddCommand(NewHiddenCmd())
	rootCmd.AddCommand(NewApplyCmd())
	// root persistentFlags

	return rootCmd
}
