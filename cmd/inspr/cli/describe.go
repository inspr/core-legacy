package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/inspr/inspr/pkg/cmd"
	cliutils "github.com/inspr/inspr/pkg/cmd/utils"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/utils"
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

	describeType := cmd.NewCmd("types <type_name | type_path>").
		WithDescription("Retrieves the full state of the type from a given namespace").
		WithExample("Display the state of the given type on the default scope", "describe types hello_world").
		WithExample("Display the state of the given type on a custom scope", "describe types --scope app1.app2 hello_world").
		WithExample("Display the state of the given type by the path", "describe types app1.app2.hello_world").
		WithAliases([]string{"t"}).
		WithCommonFlags().
		ExactArgs(1, displayTypeState)

	describeAlias := cmd.NewCmd("alias <alias_key | alais_path>").
		WithDescription("Retrieves the full state of the alias from a given namespace").
		WithExample("Display the state of the given alias on the default scope", "describe alias myalias").
		WithExample("Display the state of the given alias on a custom scope", "describe alias --scope app1.app2 myalias").
		WithExample("Display the state of the given alias by the path", "describe alias app1.app2.myalias").
		WithAliases([]string{"al"}).
		WithCommonFlags().
		ExactArgs(1, displayAlias)

	describeCmd := cmd.NewCmd("describe").
		WithDescription("Retrieves the full state of a component from a given namespace").
		WithExample("Describes the app component type", "describe apps <namespace>").
		WithExample("Describes the app component type", "describe a <namespace> --scope <specific-scope>").
		WithExample("Describes the channel component type", "describe ch <namespace>").
		WithExample("Describes the alias component type", "describe al <namespace>").
		WithLongDescription("describe takes a component type (apps | channels | types | alias) plus the name of the component, and displays the state tree)").
		AddSubCommand(describeApp).
		AddSubCommand(describeChannel).
		AddSubCommand(describeType).
		AddSubCommand(describeAlias).
		Super()

	return describeCmd
}

func displayAppState(_ context.Context, args []string) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()

	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	if !utils.IsValidScope(args[0]) {
		fmt.Fprint(out, "invalid args\n")
		return ierrors.NewError().Message("Invalid args").BadRequest().Build()
	}

	path, _ := utils.JoinScopes(scope, args[0])

	app, err := client.Apps().Get(context.Background(), path)
	if err != nil {
		cliutils.RequestErrorMessage(err, out)
		return err
	}

	utils.PrintAppTree(app, out)

	return nil
}

func displayChannelState(_ context.Context, args []string) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()

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
		cliutils.RequestErrorMessage(err, out)
		return err
	}
	utils.PrintChannelTree(channel, out)

	return nil
}

func displayTypeState(_ context.Context, args []string) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()

	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	path, typeName, err := cliutils.ProcessArg(args[0], scope)
	if err != nil {
		return err
	}

	insprType, err := client.Types().Get(context.Background(), path, typeName)
	if err != nil {
		cliutils.RequestErrorMessage(err, out)
		return err
	}
	utils.PrintTypeTree(insprType, out)

	return nil
}

func displayAlias(_ context.Context, args []string) error {
	client := cliutils.GetCliClient()
	out := cliutils.GetCliOutput()

	scope, err := cliutils.GetScope()
	if err != nil {
		return err
	}

	path, aliasKey, err := cliutils.ProcessArg(args[0], scope)
	if err != nil {
		return err
	}

	alias, err := client.Alias().Get(context.Background(), path, aliasKey)
	if err != nil {
		cliutils.RequestErrorMessage(err, out)
		return err
	}

	utils.PrintAliasTree(alias, out)

	return nil
}
