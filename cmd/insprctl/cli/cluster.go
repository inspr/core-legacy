package cli

import (
	"context"
	"fmt"

	"github.com/inspr/inspr/pkg/cmd"
	"github.com/inspr/inspr/pkg/cmd/utils"
	"github.com/spf13/cobra"
)

// NewClusterCommand creates cluster command for Inspr CLI
func NewClusterCommand() *cobra.Command {
	getBrokers := cmd.NewCmd("brokers").
		WithDescription("Retrieves brokers currently installed").
		WithExample("get cluster's brokers", "cluster brokers").
		WithAliases("b").
		NoArgs(getBrokers)
	authInit := cmd.NewCmd("init").
		WithDescription("Init configures insprd's default token").
		WithExample("init insprd as admin", "cluster init <admin_password>").
		WithCommonFlags().
		ExactArgs(1, authInit)
	return cmd.NewCmd("cluster").
		WithDescription("Configure aspects of your inspr cluster").
		WithLongDescription("Cluster takes a subcommand of (brokers | init)").
		WithExample("get cluster's brokers", "cluster brokers").
		WithExample("init insprd as admin", "cluster init <admin_password>").
		AddSubCommand(getBrokers, authInit).
		Super()
}

func getBrokers(ctx context.Context) error {
	client := utils.GetCliClient()
	out := utils.GetCliOutput()
	resp, err := client.Brokers().Get(context.Background())
	if err != nil {
		utils.RequestErrorMessage(err, out)
		return err
	}

	fmt.Fprintf(out, "DEFAULT:\n%s\n", resp.Default)
	fmt.Fprintln(out, "AVAILABLE:")
	lines := make([]string, 0)
	for _, broker := range resp.Installed {
		printLine(broker, &lines)
	}
	printTab(&lines)
	return nil
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
