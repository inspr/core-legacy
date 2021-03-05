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

// NewDescribeCmd creates describe command for Inspr CLI
func NewDescribeCmd() *cobra.Command {
	describeApp := cmd.NewCmd("apps <app_name | app_path>").
		WithDescription("Retrieves the full state of the app from a given namespace").
		WithExample("Display the state of the given app on the default scope", "describe apps hello_world").
		WithExample("Display the state of the given app on a custom scope", "describe apps --scope app1.app2 hello_world").
		WithExample("Display the state of the given app by the path", "describe apps app1.app2.hello_world").
		WithAliases([]string{"a"}).
		WithCommonFlags().
		ExactArgs(1, displayAppState)

	describeChannel := cmd.NewCmd("channels <channel_name | channel_path>").
		WithDescription("Retrieves the full state of the channel from a given namespace").
		WithExample("Display the state of the given channel on the default scope", "describe channels hello_world").
		WithExample("Display the state of the given channel on a custom scope", "describe channels --scope app1.app2 hello_world").
		WithExample("Display the state of the given channel by the path", "describe channels app1.app2.hello_world").
		WithAliases([]string{"ch"}).
		WithCommonFlags().
		ExactArgs(1, displayChannelState)

	describeChannelType := cmd.NewCmd("ctypes <ctype_name | ctype_path>").
		WithDescription("Retrieves the full state of the channelType from a given namespace").
		WithExample("Display the state of the given channelType on the default scope", "describe ctypes hello_world").
		WithExample("Display the state of the given channelType on a custom scope", "describe ctypes --scope app1.app2 hello_world").
		WithExample("Display the state of the given channelType by the path", "describe ctypes app1.app2.hello_world").
		WithAliases([]string{"ct"}).
		WithCommonFlags().
		ExactArgs(1, displayChannelTypeState)

	describeCmd := cmd.NewCmd("describe").
		WithDescription("Retrieves the full state of a component from a given namespace").
		WithExample("Describes the app component type", "describe apps <namespace>").
		WithExample("Describes the app component type", "describe a <namespace> --scope <specific-scope>").
		WithExample("Describes the channel component type", "describe ch <namespace>").
		WithExample("Describes the channel component type", "describe ch <namespace> --scope <specific-scope>").
		WithExample("Describes the channel_type component type", "describe ct <namespace>").
		WithExample("Describes the channel_type component type", "describe ct <namespace> --scope <specific-scope>").
		WithLongDescription("describe takes a component type (apps | channels | ctypes) plus the name of the component, and displays the state tree)").
		AddSubCommand(describeApp).
		AddSubCommand(describeChannel).
		AddSubCommand(describeChannelType).
		Super()

	return describeCmd
}

func displayAppState(_ context.Context, out io.Writer, args []string) error {
	client := cliutils.GetClient()

	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	if !utils.IsValidScope(args[0]) {
		fmt.Fprint(out, "invalid args\n")
		return ierrors.NewError().Message("Invalid args").BadRequest().Build()
	}

	separator := ""
	if scope != "" {
		separator = "."
	}
	path := scope + separator + args[0]

	app, err := client.Apps().Get(context.Background(), path)
	if err != nil {
		fmt.Fprintln(out, err.Error())
		return err
	}

	utils.PrintAppTree(app)

	return nil
}

func displayChannelState(_ context.Context, out io.Writer, args []string) error {
	client := cliutils.GetClient()
	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	path, chName, err := cliutils.ProcessArg(args[0], scope)
	if err != nil {
		return err
	}

	channel, err := client.Channels().Get(context.Background(), path, chName)
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}
	utils.PrintChannelTree(channel)

	return nil
}

func displayChannelTypeState(_ context.Context, out io.Writer, args []string) error {
	client := cliutils.GetClient()
	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	path, ctName, err := cliutils.ProcessArg(args[0], scope)
	if err != nil {
		return err
	}

	channelType, err := client.ChannelTypes().Get(context.Background(), path, ctName)
	if err != nil {
		fmt.Fprint(out, err.Error()+"\n")
		return err
	}
	utils.PrintChannelTypeTree(channelType)

	return nil
}
