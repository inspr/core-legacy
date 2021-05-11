package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	cliutils "github.com/inspr/inspr/pkg/cmd/utils"

	"github.com/inspr/inspr/pkg/cmd"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/utils"
)

// NewDeleteCmd creates delete command for Inspr CLI
func NewDeleteCmd() *cobra.Command {
	deleteApps := cmd.NewCmd("apps").
		WithDescription("Delete apps from context ").
		WithAliases([]string{"a"}).
		WithExample("Delete app from the default scope", "delete apps <appname> ").
		WithExample("Delete app from a custom scope", "delete apps <appname> --scope app1.app2").
		WithCommonFlags().
		MinimumArgs(1, deleteApps)
	deleteChannels := cmd.NewCmd("channels").
		WithDescription("Delete channels from context").
		WithExample("Delete channel from the default scope", "delete channels <channelname>").
		WithExample("Delete channels from a custom scope", "delete channels <channelname> --scope app1.app2").
		WithAliases([]string{"ch"}).
		WithCommonFlags().
		MinimumArgs(1, deleteChannels)
	deleteTypes := cmd.NewCmd("types").
		WithDescription("Delete types from context").
		WithExample("Delete type from the default scope", "delete types <typename>").
		WithExample("Delete type from a custom scope", "delete types <typename> --scope app1.app2").
		WithAliases([]string{"ct"}).
		WithCommonFlags().
		MinimumArgs(1, deletetypes)

	deleteAlias := cmd.NewCmd("alias").
		WithDescription("Delete alias from context").
		WithExample("Delete alias from default scope", "delete alias <aliaskey>").
		WithExample("Delete alias from a custom scope", "delete alias <aliaskey> --scope app1.app2").
		WithAliases([]string{"al"}).
		WithCommonFlags().
		MinimumArgs(1, deleteAlias)

	return cmd.NewCmd("delete").
		WithDescription("Delete component of object type").
		WithLongDescription("Delete takes a component type (apps | channels | types | alias) its scope and name, and deletes it from the cluster").
		WithExample("deletes app", "delete apps <app_name>").
		WithExample("deletes channel", "delete ch <channel_name>").
		WithExample("deletes type", "delete ct <type_name>").
		WithExample("deletes alias", "delete al <alias_key>").
		AddSubCommand(deleteApps).
		AddSubCommand(deleteChannels).
		AddSubCommand(deleteTypes).
		AddSubCommand(deleteAlias).
		Super()

}

func deleteApps(_ context.Context, args []string) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()

	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	for _, arg := range args {
		if !utils.IsValidScope(arg) {
			fmt.Fprint(out, "invalid args\n")
			return ierrors.NewError().Message("Invalid args").BadRequest().Build()
		}
		path, _ := utils.JoinScopes(scope, arg)

		cl, err := client.Apps().Delete(
			context.Background(),
			path,
			cmd.InsprOptions.DryRun,
		)
		if err != nil {
			cliutils.RequestErrorMessage(err, out)
			return err
		}
		cl.Print(out)
	}

	return nil
}

func deleteChannels(_ context.Context, args []string) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()
	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	for _, arg := range args {
		path, chName, err := cliutils.ProcessArg(arg, scope)
		if err != nil {
			return err
		}

		cl, err := client.Channels().Delete(
			context.Background(),
			path,
			chName,
			cmd.InsprOptions.DryRun,
		)
		if err != nil {
			cliutils.RequestErrorMessage(err, out)
			return err
		}
		cl.Print(out)
	}

	return nil
}

func deletetypes(_ context.Context, args []string) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()

	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	for _, arg := range args {
		path, ctName, err := cliutils.ProcessArg(arg, scope)
		if err != nil {
			return err
		}

		cl, err := client.Types().Delete(
			context.Background(),
			path,
			ctName,
			cmd.InsprOptions.DryRun,
		)
		if err != nil {
			cliutils.RequestErrorMessage(err, out)
			return err
		}
		cl.Print(out)
	}

	return nil
}

func deleteAlias(_ context.Context, args []string) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()

	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	for _, arg := range args {
		path, aliasKey, err := cliutils.ProcessArg(arg, scope)
		if err != nil {
			return err
		}

		cl, err := client.Alias().Delete(
			context.Background(),
			path,
			aliasKey,
			cmd.InsprOptions.DryRun,
		)
		if err != nil {
			cliutils.RequestErrorMessage(err, out)
			return err
		}
		cl.Print(out)
	}

	return nil
}
