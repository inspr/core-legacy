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

// NewBrokerCmd creates broker command for Inspr CLI
func NewBrokerCmd() *cobra.Command {

	kafkaCmd := cmd.NewCmd("kafka").
		WithDescription("Configures a kafka broker on insprd by importing a valid yaml file carring configurations for one of the supported brokers").
		WithExample("kafka kafka.yaml", "brokers kafka <file>").
		ExactArgs(1, kafkaConfig)

	return cmd.NewCmd("brokers").
		WithDescription("Retrieves brokers currently installed").
		WithLongDescription(`Broker is the command that returns the brokers already installed on the cluster,
		and has the kafka subcommand wich installs the kafka broker on the cluster`).
		WithExample("get brokers", "brokers").
		WithExample("install kafka broker from a kafka.yaml", "brokers kafka <file>").
		AddSubCommand(kafkaCmd).
		NoArgs(getBrokers)
}

func getBrokers(ctx context.Context) error {
	client := utils.GetCliClient()
	out := utils.GetCliOutput()
	resp, err := client.Brokers().Get(context.Background())
	if err != nil {
		fmt.Fprint(out, ierrors.FormatError(err))
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

func kafkaConfig(c context.Context, args []string) error {
	return brokerConfig("kafka", args[0])
}

func brokerConfig(brokerName, filePath string) error {
	client := utils.GetCliClient()
	output := utils.GetCliOutput()

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
		return ierrors.New("not a yaml file").InvalidFile()
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
