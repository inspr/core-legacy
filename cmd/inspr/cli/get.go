package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/spf13/cobra"
	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

var tabWriter *tabwriter.Writer
var lines []string

// NewGetCmd creates get command for Inspr CLI
func NewGetCmd() *cobra.Command {
	getApps := cmd.NewCmd("apps").
		WithDescription("Get apps from context ").
		WithAliases([]string{"a"}).
		WithExample("Get apps from the default scope", "get apps ").
		WithExample("Get apps from a custom scope", "get apps --scope app1.app2").
		WithCommonFlags().
		NoArgs(getApps)
	getChannels := cmd.NewCmd("channels").
		WithDescription("Get channels from context").
		WithExample("Get channels from the default scope", "get channels ").
		WithExample("Get channels from a custom scope", "get channels --scope app1.app2").
		WithAliases([]string{"ch"}).
		WithCommonFlags().
		NoArgs(getChannels)
	getTypes := cmd.NewCmd("ctypes").
		WithDescription("Get channel types from context").
		WithExample("Get channel types from the default scope", "get ctypes ").
		WithExample("Get channel types from a custom scope", "get ctypes --scope app1.app2").
		WithAliases([]string{"ct"}).
		WithCommonFlags().
		NoArgs(getCTypes)
	getNodes := cmd.NewCmd("nodes").
		WithDescription("Get nodes from context").
		WithExample("Get nodes from the default scope", "get nodes ").
		WithExample("Get nodes from a custom scope", "get nodes --scope app1.app2").
		WithAliases([]string{"n"}).
		WithCommonFlags().
		NoArgs(getNodes)
	return cmd.NewCmd("get").
		WithDescription("Get by object type").
		WithDescription("Retrieves the components from a given namespace").
		WithLongDescription("get takes a component type (apps | channels | ctypes | nodes) and displays names for those components is a scope)").
		WithAliases([]string{"list"}).
		AddSubCommand(getApps).
		AddSubCommand(getChannels).
		AddSubCommand(getTypes).
		AddSubCommand(getNodes).
		Super()

}

func getApps(_ context.Context, out io.Writer) error {
	initTab(out)
	err := getObj(printApps, out)
	if err != nil {
		return err
	}
	printTab()
	return nil
}

func getChannels(_ context.Context, out io.Writer) error {
	initTab(out)
	err := getObj(printChannels, out)
	if err != nil {
		return err
	}
	printTab()
	return nil
}

func getCTypes(_ context.Context, out io.Writer) error {
	initTab(out)
	err := getObj(printCTypes, out)
	if err != nil {
		return err
	}
	printTab()
	return nil
}

func getNodes(_ context.Context, out io.Writer) error {
	initTab(out)
	err := getObj(printNodes, out)
	if err != nil {
		return err
	}
	printTab()
	return nil
}

func getObj(printObj func(*meta.App), out io.Writer) error {
	rc := request.NewClient().BaseURL(cliutils.GetConfiguredServerIP()).Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build()
	client := client.NewControllerClient(rc)
	scope, err := cliutils.GetScope()
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}
	resp, err := client.Apps().Get(context.Background(), scope)
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}
	printObj(resp)
	return nil
}

func printApps(app *meta.App) {
	if app.Meta.Name != "" {
		printLine(app.Meta.Name)
	}
	for _, child := range app.Spec.Apps {
		printApps(child)
	}
}

func printChannels(app *meta.App) {
	for ch := range app.Spec.Channels {
		printLine(ch)
	}
	for _, child := range app.Spec.Apps {
		printChannels(child)
	}
}

func printCTypes(app *meta.App) {
	for ct := range app.Spec.ChannelTypes {
		printLine(ct)
	}
	for _, child := range app.Spec.Apps {
		printChannels(child)
	}
}

func printNodes(app *meta.App) {
	if app.Spec.Node.Meta.Name != "" {
		printLine(app.Spec.Node.Meta.Name)
	}
	for _, child := range app.Spec.Apps {
		printApps(child)
	}
}

func printLine(name string) {
	lines = append(lines, fmt.Sprintf("%s\n", name))
}

func initTab(out io.Writer) {
	tabWriter = tabwriter.NewWriter(out, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	lines = append(lines, "NAME\n")
}

func printTab() {
	for _, line := range lines {
		fmt.Fprint(tabWriter, line)
	}
	tabWriter.Flush()
}
