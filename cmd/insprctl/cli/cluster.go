package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"inspr.dev/inspr/pkg/cmd"
	"inspr.dev/inspr/pkg/cmd/utils"
	"inspr.dev/inspr/pkg/ierrors"
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
	configCmd := cmd.NewCmd("config").
		WithDescription("configures a broker on insprd by importing a valid yaml file carring configurations for one of the supported brokers").
		WithExample("config kafka kafka.yaml", "cluster config <broker> <file>").
		ExactArgs(2, clusterConfig)
	return cmd.NewCmd("cluster").
		WithDescription("Configure aspects of your insprd cluster").
		WithLongDescription("Cluster takes a subcommand of (brokers | init)").
		WithExample("get cluster's brokers", "cluster brokers").
		WithExample("init insprd as admin", "cluster init <admin_password>").
		AddSubCommand(getBrokers, authInit, configCmd).
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
	for _, broker := range resp.Available {
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

func clusterConfig(c context.Context, args []string) error {
	client := utils.GetCliClient()
	output := utils.GetCliOutput()
	brokerName, filePath := args[0], args[1]

	if err := utils.CheckEmptyArgs(map[string]string{
		"brokerName": brokerName,
		"filePath":   filePath,
	}); err != nil {
		fmt.Fprintf(output, "invalid args: %v\n", err.Error())
		return err
	}

	// check if file exists and if it is a yaml file
	if _, err := os.Stat(filePath); os.IsNotExist(err) || !isYaml(filePath) {
		if err != nil {
			fmt.Fprintf(output, "unable to find file: %v\n", err.Error())
			return err
		}

		fmt.Fprintf(output, "not a yaml file\n")
		return ierrors.NewError().Message("not a yaml file").InvalidFile().Build()
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(output, "unable to read file: %v\n", err.Error())
		return err
	}

	// do a request to the broker route /brokers/<broker_name>
	err = client.Brokers().Create(context.Background(), brokerName, bytes)
	if err != nil {
		fmt.Fprintf(output, "unable to create broker: %v\n", err.Error())
		return err
	}

	fmt.Fprintln(output, "successfully installed broker on insprd")
	return nil
}
