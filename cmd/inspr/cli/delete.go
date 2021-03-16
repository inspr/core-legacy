package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"

	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
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
	deleteTypes := cmd.NewCmd("ctypes").
		WithDescription("Delete channel types from context").
		WithExample("Delete channel type from the default scope", "delete ctypes <ctypename>").
		WithExample("Delete channel type from a custom scope", "delete ctypes <ctypename> --scope app1.app2 ").
		WithAliases([]string{"ct"}).
		WithCommonFlags().
		MinimumArgs(1, deleteCTypes)
	return cmd.NewCmd("delete").
		WithDescription("Delete component of object type").
		WithLongDescription("Delete takes a component type (apps | channels | ctypes) its scope and name, and deletes it from the cluster").
		WithExample("deletes app", "delete apps <app_name>").
		WithExample("deletes channel", "delete ch <channel_name>").
		WithExample("deletes channel_type", "delete ct <channel_type_name>").
		AddSubCommand(deleteApps).
		AddSubCommand(deleteChannels).
		AddSubCommand(deleteTypes).
		Super()

}

func deleteApps(_ context.Context, args []string) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()

	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	for idx := range args {
		if !utils.IsValidScope(args[idx]) {
			fmt.Fprint(out, "invalid args\n")
			return ierrors.NewError().Message("Invalid args").BadRequest().Build()
		}
		path, _ := utils.JoinScopes(scope, args[idx])

		cl, err := client.Apps().Delete(context.Background(), path, cmd.InsprOptions.DryRun)
		if err != nil {
			fmt.Fprint(out, err.Error()+"\n")
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

	for idx := range args {
		path, chName, err := cliutils.ProcessArg(args[idx], scope)
		if err != nil {
			return err
		}

		cl, err := client.Channels().Delete(context.Background(), path, chName, cmd.InsprOptions.DryRun)
		if err != nil {
			fmt.Fprint(out, err.Error()+"\n")
			return err
		}
		cl.Print(out)
	}

	return nil
}

func deleteCTypes(_ context.Context, args []string) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()

	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	for idx := range args {
		path, ctName, err := cliutils.ProcessArg(args[idx], scope)
		if err != nil {
			return err
		}

		cl, err := client.ChannelTypes().Delete(context.Background(), path, ctName, cmd.InsprOptions.DryRun)
		if err != nil {
			fmt.Fprint(out, err.Error()+"\n")
			return err
		}
		cl.Print(out)
	}

	return nil
}
