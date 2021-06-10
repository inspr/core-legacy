package cli

import (
	"context"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/spf13/cobra"

	cliutils "github.com/inspr/inspr/pkg/cmd/utils"
	"github.com/inspr/inspr/pkg/controller"
	"github.com/inspr/inspr/pkg/ierrors"

	"github.com/inspr/inspr/pkg/cmd"
	"github.com/inspr/inspr/pkg/meta"
)

// NewGetCmd creates get command for Inspr CLI
func NewGetCmd() *cobra.Command {
	getApps := cmd.NewCmd("apps").
		WithDescription("Get apps from context ").
		WithAliases("a").
		WithExample("Get apps from the default scope", "get apps ").
		WithExample("Get apps from a custom scope", "get apps --scope app1.app2").
		WithCommonFlags().
		NoArgs(getApps)
	getChannels := cmd.NewCmd("channels").
		WithDescription("Get channels from context").
		WithExample("Get channels from the default scope", "get channels ").
		WithExample("Get channels from a custom scope", "get channels --scope app1.app2").
		WithAliases("ch").
		WithCommonFlags().
		NoArgs(getChannels)
	getTypes := cmd.NewCmd("types").
		WithDescription("Get types from context").
		WithExample("Get types from the default scope", "get types ").
		WithExample("Get types from a custom scope", "get types --scope app1.app2").
		WithAliases("t").
		WithCommonFlags().
		NoArgs(getTypes)
	getNodes := cmd.NewCmd("nodes").
		WithDescription("Get nodes from context").
		WithExample("Get nodes from the default scope", "get nodes ").
		WithExample("Get nodes from a custom scope", "get nodes --scope app1.app2").
		WithAliases("n").
		WithCommonFlags().
		NoArgs(getNodes)
	getAlias := cmd.NewCmd("alias").
		WithDescription("Get alias from context").
		WithExample("Get alias from the default scope", "get alias ").
		WithExample("Get alias from a custom scope", "get alias --scope app1.app2").
		WithAliases("al").
		WithCommonFlags().
		NoArgs(getAlias)
	return cmd.NewCmd("get").
		WithDescription("Get by object type").
		WithDescription("Retrieves the components from a given namespace").
		WithExample("gets apps from cluster", "get apps --scope <scope>").
		WithExample("gets channels from cluster", "get ch --scope <scope>").
		WithExample("gets types from cluster", "get t --scope <scope>").
		WithExample("gets nodes from cluster", "get nodes --scope <scope>").
		WithExample("gets alias from cluster", "get alias --scope <scope>").
		WithLongDescription("get takes a component type (apps | channels | types | nodes | alias) and displays names for those components is a scope)").
		WithAliases("list").
		AddSubCommand(getApps).
		AddSubCommand(getChannels).
		AddSubCommand(getTypes).
		AddSubCommand(getNodes).
		AddSubCommand(getAlias).
		Super()

}

func getApps(_ context.Context) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()
	scope, err := cliutils.GetScope()
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}

	lines := make([]string, 0)
	initTab(&lines)
	err = getObj(printApps, &lines, client, out, scope)
	if err != nil {
		return err
	}
	printTab(&lines)
	return nil
}

func getChannels(_ context.Context) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()
	scope, err := cliutils.GetScope()
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}

	_, err = client.Channels().Get(context.Background(), scope, "")
	if ierrors.HasCode(err, ierrors.Forbidden) {
		cliutils.RequestErrorMessage(err, out)
		return err
	}

	lines := make([]string, 0)
	initTab(&lines)

	err = getObj(printChannels, &lines, client, out, scope)
	if err != nil {
		return err
	}
	printTab(&lines)
	return nil
}

func getTypes(_ context.Context) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()
	scope, err := cliutils.GetScope()
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}

	_, err = client.Types().Get(context.Background(), scope, "")
	if ierrors.HasCode(err, ierrors.Forbidden) {
		cliutils.RequestErrorMessage(err, out)
		return err
	}

	lines := make([]string, 0)
	initTab(&lines)

	err = getObj(printTypes, &lines, client, out, scope)
	if err != nil {
		return err
	}

	printTab(&lines)
	return nil
}

func getAlias(_ context.Context) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()
	scope, err := cliutils.GetScope()
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}

	_, err = client.Alias().Get(context.Background(), scope, "")
	if ierrors.HasCode(err, ierrors.Forbidden) {
		cliutils.RequestErrorMessage(err, out)
		return err
	}

	lines := make([]string, 0)
	initTab(&lines)

	err = getObj(printAliases, &lines, client, out, scope)
	if err != nil {
		return err
	}

	printTab(&lines)
	return nil
}

func getNodes(_ context.Context) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()
	scope, err := cliutils.GetScope()
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}

	lines := make([]string, 0)
	initTab(&lines)

	err = getObj(printNodes, &lines, client, out, scope)
	if err != nil {
		return err
	}

	printTab(&lines)
	return nil
}

func getObj(printObj func(*meta.App, *[]string), lines *[]string, client controller.Interface, out io.Writer, scope string) error {
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

func printTypes(app *meta.App, lines *[]string) {
	for insprType := range app.Spec.Types {
		printLine(insprType, lines)
	}
	for _, child := range app.Spec.Apps {
		printTypes(child, lines)
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
