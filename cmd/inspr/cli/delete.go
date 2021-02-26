package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
)

// NewDeleteCmd - mock subcommand
func NewDeleteCmd() *cobra.Command {
	deleteApps := cmd.NewCmd("apps").
		WithDescription("Delete apps from context ").
		WithAliases([]string{"a"}).
		WithExample("Delete apps from the default scope", "delete apps ").
		WithExample("Delete apps from a custom scope", "delete apps --scope app1.app2").
		WithCommonFlags().
		ExactArgs(1, deleteApps)
	deleteChannels := cmd.NewCmd("channels").
		WithDescription("Delete channels from context").
		WithExample("Delete channels from the default scope", "delete channels ").
		WithExample("Delete channels from a custom scope", "delete channels --scope app1.app2").
		WithAliases([]string{"ch"}).
		WithCommonFlags().
		ExactArgs(1, deleteChannels)
	deleteTypes := cmd.NewCmd("ctypes").
		WithDescription("Delete channel types from context").
		WithExample("Delete channel types from the default scope", "delete ctypes ").
		WithExample("Delete channel types from a custom scope", "delete ctypes --scope app1.app2").
		WithAliases([]string{"ct"}).
		WithCommonFlags().
		ExactArgs(1, deleteCTypes)
	return cmd.NewCmd("delete").
		WithDescription("Delete by object type").
		WithDescription("Retrieves the components from a given namespace").
		WithLongDescription("Delete takes a component type (apps | channels | ctypes) its scope and name, and deletes it from the cluster").
		AddSubCommand(deleteApps).
		AddSubCommand(deleteChannels).
		AddSubCommand(deleteTypes).
		Super()

}

func deleteApps(_ context.Context, out io.Writer, args []string) error {
	client := getClient()

	scope, err := getScope()
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
	client := getClient()
	scope, err := getScope()
	if err != nil {
		return err
	}

	path, chName, err := processArg(args[0], scope)
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
	client := getClient()
	scope, err := getScope()
	if err != nil {
		return err
	}

	path, ctName, err := processArg(args[0], scope)
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
