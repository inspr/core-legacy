package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

var tabWriter *tabwriter.Writer
var ctx string

// NewGetCmd - mock subcommand
func NewGetCmd() *cobra.Command {
	getApps := cmd.NewCmd("apps").
		WithDescription("Get apps").
		WithAliases([]string{"a"}).
		WithCommonFlags().
		WithFlags([]*cmd.Flag{
			{
				Name:          "scope",
				Usage:         "inspr get <subcommand> --scope/-s <apppath>",
				Shorthand:     "s",
				Value:         &ctx,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"apps"},
			},
		}).
		NoArgs(getApps)
	getChannels := cmd.NewCmd("channels").
		WithDescription("Get channels").
		WithAliases([]string{"ch"}).
		WithCommonFlags().
		WithFlags([]*cmd.Flag{
			{
				Name:          "define search scope",
				Usage:         "inspr get <subcommand> --scope/-s <apppath>",
				Shorthand:     "s",
				Value:         &ctx,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"channels"},
			},
		}).
		NoArgs(getChannels)
	getTypes := cmd.NewCmd("types").
		WithDescription("Get types").
		WithAliases([]string{"ct"}).
		WithCommonFlags().
		WithFlags([]*cmd.Flag{
			{
				Name:          "scope",
				Usage:         "inspr get <subcommand> --scope/-s <apppath>",
				Shorthand:     "s",
				Value:         &ctx,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"types"},
			},
		}).
		NoArgs(getCTypes)
	getNodes := cmd.NewCmd("nodes").
		WithDescription("Get nodes").
		WithAliases([]string{"n"}).
		WithCommonFlags().
		WithFlags([]*cmd.Flag{
			{
				Name:          "scope",
				Usage:         "inspr get <subcommand> --scope/-s <apppath>",
				Shorthand:     "s",
				Value:         &ctx,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"nodes"},
			},
		}).
		NoArgs(getNodes)
	return cmd.NewCmd("get").
		WithDescription("Get by object type").
		WithAliases([]string{"list"}).
		WithCommonFlags().
		AddSubCommand(getApps).
		AddSubCommand(getChannels).
		AddSubCommand(getTypes).
		AddSubCommand(getNodes).
		Super()

}

func getApps(_ context.Context, out io.Writer) error {
	initTab(out)
	getObj(printApps)
	printTab()
	return nil
}

func getChannels(_ context.Context, out io.Writer) error {
	initTab(out)
	getObj(printChannels)
	printTab()
	return nil
}

func getCTypes(_ context.Context, out io.Writer) error {
	initTab(out)
	getObj(printCTypes)
	printTab()
	return nil
}

func getNodes(_ context.Context, out io.Writer) error {
	initTab(out)
	getObj(printNodes)
	printTab()
	return nil
}

func getObj(printObj func(*meta.App)) {
	getDO := models.AppQueryDI{
		Ctx:    ctx,
		Valid:  true,
		DryRun: false,
	}
	body, err := json.Marshal(getDO)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req, err := http.NewRequest(http.MethodGet, getAppsURL(), bytes.NewBuffer(body))
	defer req.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	apps := &models.AppDI{}
	json.NewDecoder(resp.Body).Decode(apps)
	printObj(&apps.App)

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
	fmt.Fprintf(tabWriter, "%s", name)
}

func initTab(out io.Writer) {
	tabWriter = tabwriter.NewWriter(out, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintf(tabWriter, "NAME")
}

func printTab() {
	tabWriter.Flush()
}

func getAppsURL() string {
	return fmt.Sprintf(viper.GetString("reqUrl"), "/apps")
}
