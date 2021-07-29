package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"inspr.dev/inspr/pkg/cmd"
	"inspr.dev/inspr/pkg/cmd/utils"
)

// NewClusterCommand creates cluster command for Inspr CLI
func NewClusterCommand() *cobra.Command {
	authInit := cmd.NewCmd("init").
		WithDescription("Init configures insprd's default token").
		WithExample("init insprd as admin", "cluster init <admin_password>").
		WithCommonFlags().
		ExactArgs(1, authInit)
	return cmd.NewCmd("cluster").
		WithDescription("Configure aspects of your insprd cluster").
		WithLongDescription("Cluster takes a subcommand of (init)").
		WithExample("init insprd as admin", "cluster init <admin_password>").
		AddSubCommand(authInit).
		Super()
}

func authInit(c context.Context, args []string) error {
	output := utils.GetCliOutput()

	token, err := utils.GetCliClient().Authorization().Init(c, args[0])
	if err != nil {
		utils.RequestErrorMessage(err, output)
		return err
	}

	fmt.Fprintln(output, "This is a root token for authentication within your insprd. This will not be generated again. Save it wisely.")
	fmt.Fprintf(output, "%s\n", token)
	return nil
}
