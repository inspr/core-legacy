package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/inspr/inspr/pkg/cmd"
	"github.com/inspr/inspr/pkg/cmd/utils"
	cliutils "github.com/inspr/inspr/pkg/cmd/utils"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/spf13/cobra"
)

// NewClusterCommand creates cluster command for Inspr CLI
func NewClusterCommand() *cobra.Command {
	getBrokers := cmd.NewCmd("brokers").
		WithDescription("Retrieves brokers currently installed").
		WithExample("get cluster's brokers", "inspr cluster brokers").
		WithAliases("b").
		NoArgs(getBrokers)
	authInit := cmd.NewCmd("init").
		WithDescription("Init configures insprd's default token").
		WithExample("init insprd as admin", " inspr cluster init <admin_password>").
		WithCommonFlags().
		ExactArgs(1, authInit)
	configCmd := cmd.NewCmd("config").
		WithDescription("obtains the broker and yaml file and tries to install it on the insprd server").
		WithExample("config kafka kafka.yaml", "inspr cluster config <broker> <file>").
		WithCommonFlags().
		ExactArgs(2, clusterConfig)
	return cmd.NewCmd("cluster").
		WithDescription("Configure aspects of your inspr cluster").
		WithLongDescription("Cluster takes a subcommand of (brokers | init)").
		WithExample("get cluster's brokers", "inspr cluster brokers").
		WithExample("init insprd as admin", " inspr cluster init <admin_password>").
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

func clusterConfig(c context.Context, args []string) error {
	client := cliutils.GetCliClient()
	output := utils.GetCliOutput()
	brokerName, filePath := args[0], args[1]
	// check if file exists and if it is a yaml file
	if _, err := os.Stat(filePath); os.IsNotExist(err) || !isYaml(filePath) {
		fmt.Fprintln(output, err)
		return err
	}
	// check if broker arg is empty
	// TODO is it necessary to check? if empty it isn't a arg
	if brokerName == "" {
		fmt.Fprintln(output, "")
		return errors.New("empty brokerName")
	}

	// how to get the config according to the broker_name ?
	fileContent, _ := os.ReadFile(filePath)
	brokerData, err := metautils.YamlToKafkaConfig(fileContent)
	if err != nil {
		return err
	}

	// do a request to the kafka route /brokers/<broker_name>
	req, _ := http.NewRequest(http.MethodPost, "brokers/"+brokerName, bytes.NewBuffer([]byte(brokerName)))
	Json

	// return error message or show it was successful

	return nil
}
