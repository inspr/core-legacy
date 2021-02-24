package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

var deleteScope string

// NewGetCmd - mock subcommand
func NewDeleteCmd() *cobra.Command {
	getApps := cmd.NewCmd("apps").
		WithDescription("Get apps").
		WithAliases([]string{"a"}).
		WithCommonFlags().
		NoArgs(getApps)
	getChannels := cmd.NewCmd("channels").
		WithDescription("Get channels").
		WithAliases([]string{"ch"}).
		WithCommonFlags().
		NoArgs(getChannels)
	getTypes := cmd.NewCmd("types").
		WithDescription("Get types").
		WithAliases([]string{"ct"}).
		WithCommonFlags().
		NoArgs(getCTypes)
	getNodes := cmd.NewCmd("nodes").
		WithDescription("Get nodes").
		WithAliases([]string{"n"}).
		WithCommonFlags().
		NoArgs(getNodes)
	return cmd.NewCmd("get").
		WithDescription("Get by object type").
		WithAliases([]string{"list"}).
		WithCommonFlags().
		WithFlags([]*cmd.Flag{
			{
				Name:          "inspr get <subcommand> --scope/-s <apppath>",
				Usage:         "define search scope",
				Shorthand:     "s",
				Value:         &getScope,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"get"},
			},
		}).
		AddSubCommand(getApps).
		AddSubCommand(getChannels).
		AddSubCommand(getTypes).
		AddSubCommand(getNodes).
		Super()

}

func deleteApps(_ context.Context, out io.Writer) error {
	initTab(out)
	getObj(printApps)
	printTab()
	return nil
}

func deleteChannels(_ context.Context, out io.Writer) error {
	initTab(out)
	getObj(printChannels)
	printTab()
	return nil
}

func deleteCTypes(_ context.Context, out io.Writer) error {
	initTab(out)
	getObj(printCTypes)
	printTab()
	return nil
}

func deleteNodes(_ context.Context, out io.Writer) error {
	initTab(out)
	getObj(printNodes)
	printTab()
	return nil
}

func deleteObj(printObj func(*meta.App)) {
	getDO := models.AppQueryDI{
		Ctx:    getScope,
		Valid:  true,
		DryRun: false,
	}
	body, err := json.Marshal(getDO)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req, err := http.NewRequest(http.MethodGet, getURL(), bytes.NewBuffer(body))
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
