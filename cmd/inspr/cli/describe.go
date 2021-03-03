package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

// NewDescribeCmd DOC TODO
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
		WithExample("Describes the channel component type", "describe ch <namespace>").
		WithExample("Describes the channel_type component type", "describe ct <namespace>").
		WithLongDescription("describe takes a component type (apps | channels | ctypes) plus the name of the component, and displays the state tree)").
		AddSubCommand(describeApp).
		AddSubCommand(describeChannel).
		AddSubCommand(describeChannelType).
		Super()

	return describeCmd
}

func displayAppState(_ context.Context, out io.Writer, args []string) error {
	client := getClient()

	scope, err := getScope()
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
		fmt.Fprintln(out, path)
		fmt.Fprintln(out, scope)
		fmt.Fprintln(out, separator)
		fmt.Fprintln(out, args[0])
		return err
	}

	utils.PrintAppTree(app)

	return nil
}

func displayChannelState(_ context.Context, out io.Writer, args []string) error {
	client := getClient()
	scope, err := getScope()
	if err != nil {
		return err
	}

	path, chName, err := processArg(args[0], scope)
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
	client := getClient()
	scope, err := getScope()
	if err != nil {
		return err
	}

	path, ctName, err := processArg(args[0], scope)
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

func getClient() *client.Client {
	url := viper.GetString(configServerIP)

	rc := request.NewClient().BaseURL(url).Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build()
	return client.NewControllerClient(rc)
}

func getScope() (string, error) {
	defaultScope := viper.GetString(configScope)
	scope := defaultScope

	if cmd.InsprOptions.Scope != "" {
		if utils.IsValidScope(cmd.InsprOptions.Scope) {
			scope = cmd.InsprOptions.Scope
		} else {
			return "", ierrors.NewError().BadRequest().Message("invalid scope").Build()
		}
	}

	return scope, nil
}

func processArg(arg, scope string) (string, string, error) {
	path := scope
	var component string

	if err := utils.StructureNameIsValid(arg); err != nil {
		if !utils.IsValidScope(arg) {
			return "", "", ierrors.NewError().Message("invalid scope").BadRequest().Build()
		}

		newScope, lastName, err := utils.RemoveLastPartInScope(arg)
		if err != nil {
			return "", "", ierrors.NewError().Message("invalid scope").BadRequest().Build()
		}

		separator := ""
		if scope != "" {
			separator = "."
		}

		path = path + separator + newScope
		component = lastName

	} else {
		component = arg
	}
	return path, component, nil
}
