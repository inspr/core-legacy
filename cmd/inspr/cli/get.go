package cli

import (
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	cliutils "github.com/inspr/inspr/pkg/cmd/utils"

	"github.com/inspr/inspr/pkg/cmd"
	"github.com/inspr/inspr/pkg/meta"
)

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
	getAlias := cmd.NewCmd("alias").
		WithDescription("Get alias from context").
		WithExample("Get alias from the default scope", "get alias ").
		WithExample("Get alias from a custom scope", "get alias --scope app1.app2").
		WithAliases([]string{"al"}).
		WithCommonFlags().
		NoArgs(getAlias)
	return cmd.NewCmd("get").
		WithDescription("Get by object type").
		WithDescription("Retrieves the components from a given namespace").
		WithExample("gets apps from cluster", "get apps --scope <scope>").
		WithExample("gets channels from cluster", "get ch --scope <scope>").
		WithExample("gets channel_types from cluster", "get ct --scope <scope>").
		WithExample("gets nodes from cluster", "get nodes --scope <scope>").
		WithExample("gets alias from cluster", "get alias --scope <scope>").
		WithLongDescription("get takes a component type (apps | channels | ctypes | nodes | alias) and displays names for those components is a scope)").
		WithAliases([]string{"list"}).
		AddSubCommand(getApps).
		AddSubCommand(getChannels).
		AddSubCommand(getTypes).
		AddSubCommand(getNodes).
		AddSubCommand(getAlias).
		Super()

}

func getApps(_ context.Context) error {
	lines := make([]string, 0)
	initTab(&lines)
	err := getObj(printApps, &lines)
	if err != nil {
		return err
	}
	printTab(&lines)
	return nil
}

func getChannels(_ context.Context) error {
	lines := make([]string, 0)
	initTab(&lines)
	err := getObj(printChannels, &lines)
	if err != nil {
		return err
	}
	printTab(&lines)
	return nil
}

func getCTypes(_ context.Context) error {
	lines := make([]string, 0)
	initTab(&lines)
	err := getObj(printCTypes, &lines)
	if err != nil {
		return err
	}
	printTab(&lines)
	return nil
}

func getAlias(_ context.Context) error {
	lines := make([]string, 0)
	initTab(&lines)
	err := getObj(printAliases, &lines)
	if err != nil {
		return err
	}
	printTab(&lines)
	return nil
}

func getNodes(_ context.Context) error {
	lines := make([]string, 0)
	initTab(&lines)
	err := getObj(printNodes, &lines)
	if err != nil {
		return err
	}
	printTab(&lines)
	return nil
}

func getObj(printObj func(*meta.App, *[]string), lines *[]string) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()

	scope, err := cliutils.GetScope()
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}

	resp, err := client.Apps().Get(context.Background(), scope)
	if err != nil {
		cliutils.RequestErrorMessage(err, out)
		return err
	}

	printObj(resp, lines)
	return nil
}

func printApps(app *meta.App, lines *[]string) {
	if app.Meta.Name != "" {
		printLine(app.Meta.Name, lines)
	}
	for _, child := range app.Spec.Apps {
		printApps(child, lines)
	}
}

func printChannels(app *meta.App, lines *[]string) {
	for ch := range app.Spec.Channels {
		printLine(ch, lines)
	}
	for _, child := range app.Spec.Apps {
		printChannels(child, lines)
	}
}

func printCTypes(app *meta.App, lines *[]string) {
	for ct := range app.Spec.ChannelTypes {
		printLine(ct, lines)
	}
	for _, child := range app.Spec.Apps {
		printCTypes(child, lines)
	}
}

func printNodes(app *meta.App, lines *[]string) {
	if app.Spec.Node.Meta.Name != "" {
		printLine(app.Spec.Node.Meta.Name, lines)
	}
	for _, child := range app.Spec.Apps {
		printNodes(child, lines)
	}
}

func printAliases(app *meta.App, lines *[]string) {
	for alias := range app.Spec.Aliases {
		printLine(alias, lines)
	}
	for _, child := range app.Spec.Apps {
		printAliases(child, lines)
	}
}

func printLine(name string, lines *[]string) {
	*lines = append(*lines, fmt.Sprintf("%s\n", name))
}

func initTab(lines *[]string) {
	*lines = append(*lines, "NAME\n")
}

func printTab(lines *[]string) {
	out := cliutils.GetCliOutput()
	tabWriter := tabwriter.NewWriter(out, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for _, line := range *lines {
		fmt.Fprint(tabWriter, line)
	}
	tabWriter.Flush()
}
