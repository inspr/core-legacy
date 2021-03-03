package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
)

// NewDeleteCmd - mock subcommand
func NewDeleteCmd() *cobra.Command {
	deleteApps := cmd.NewCmd("apps").
		WithDescription("Delete apps from context ").
		WithAliases([]string{"a"}).
		WithExample("Delete app from the default scope", "delete apps <appname> ").
		WithExample("Delete app from a custom scope", "delete apps <appname> --scope app1.app2").
		WithCommonFlags().
		ExactArgs(1, deleteApps)
	deleteChannels := cmd.NewCmd("channels").
		WithDescription("Delete channels from context").
		WithExample("Delete channel from the default scope", "delete channels <channelname>").
		WithExample("Delete channels from a custom scope", "delete channels <channelname> --scope app1.app2").
		WithAliases([]string{"ch"}).
		WithCommonFlags().
		ExactArgs(1, deleteChannels)
	deleteTypes := cmd.NewCmd("ctypes").
		WithDescription("Delete channel types from context").
		WithExample("Delete channel type from the default scope", "delete ctypes <ctypename>").
		WithExample("Delete channel type from a custom scope", "delete ctypes <ctypename> --scope app1.app2 ").
		WithAliases([]string{"ct"}).
		WithCommonFlags().
		ExactArgs(1, deleteCTypes)
	return cmd.NewCmd("delete").
		WithDescription("Delete component of object type").
		WithLongDescription("Delete takes a component type (apps | channels | ctypes) its scope and name, and deletes it from the cluster").
		AddSubCommand(deleteApps).
		AddSubCommand(deleteChannels).
		AddSubCommand(deleteTypes).
		Super()

}

func deleteApps(_ context.Context, out io.Writer, args []string) error {
	client := cliutils.GetClient()

	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	if !utils.IsValidScope(args[0]) {
		fmt.Fprint(out, "invalid args\n")
		return ierrors.NewError().Message("Invalid args").BadRequest().Build()
	}

	path, err := utils.JoinScopes(scope, args[0])
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}

	cl, err := client.Apps().Delete(context.Background(), path, cmd.InsprOptions.DryRun)
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}
	cl.Print(out)
	return nil
}

func deleteChannels(_ context.Context, out io.Writer, args []string) error {
	client := cliutils.GetClient()
	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	path, chName, err := cliutils.ProcessArg(args[0], scope)
	if err != nil {
		return err
	}

	cl, err := client.Channels().Delete(context.Background(), path, chName, cmd.InsprOptions.DryRun)
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}
	cl.Print(out)

	return nil
}

func deleteCTypes(_ context.Context, out io.Writer, args []string) error {
	client := cliutils.GetClient()
	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	path, ctName, err := cliutils.ProcessArg(args[0], scope)
	if err != nil {
		return err
	}

	cl, err := client.ChannelTypes().Delete(context.Background(), path, ctName, cmd.InsprOptions.DryRun)
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}
	cl.Print(out)

	return nil
}
