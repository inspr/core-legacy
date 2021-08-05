package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"inspr.dev/inspr/pkg/cmd"
	"inspr.dev/inspr/pkg/cmd/utils"
	"inspr.dev/inspr/pkg/ierrors"
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
	out := utils.GetCliOutput()

	token, err := utils.GetCliClient().Authorization().Init(c, args[0])
	if err != nil {
		fmt.Fprintf(out, "%v\n", ierrors.FormatError(err))
		return err
	}

	fmt.Fprintln(out, "This is a root token for authentication within your insprd. This will not be generated again. Save it wisely.")
	fmt.Fprintf(out, "%s\n", token)
	return nil
}
