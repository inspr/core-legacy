package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

// NewDescribeCmd DOC TODO
func NewDescribeCmd() *cobra.Command {
	describeCmd := cmd.NewCmd("describe").
		WithDescription("retrieves the full state of a component from a given namespace").
		WithLongDescription("describe takes a component type (app | channel | channelType | node) plus the name of the component, and displays the state tree)").
		Super()

	describeApp := cmd.NewCmd("apps").
		WithDescription("retrieves the full state of the app from a given namespace").
		WithExample("Display the state of the given app on the default scope", "describe apps hello_world").
		WithExample("Display the state of the given app on a custom scope", "describe apps --scope app1.app2 hello_world").
		WithExample("Display the state of the given app by the path", "describe apps app1.app2.hello_world").
		WithAliases([]string{"a"}).
		WithCommonFlags().
		ExactArgs(1, displayAppState)

	describeChannel := cmd.NewCmd("channels").
		WithDescription("retrieves the full state of the channel from a given namespace").
		WithExample("Display the state of the given channel on the default scope", "describe channels hello_world").
		WithExample("Display the state of the given channel on a custom scope", "describe channels --scope app1.app2 hello_world").
		WithExample("Display the state of the given channel by the path", "describe channels app1.app2.hello_world").
		WithAliases([]string{"ch"}).
		WithCommonFlags().
		ExactArgs(1, displayChannelState)

	describeChannelType := cmd.NewCmd("channeltypes").
		WithDescription("retrieves the full state of the channelType from a given namespace").
		WithExample("Display the state of the given channelType on the default scope", "describe channeltypes hello_world").
		WithExample("Display the state of the given channelType on a custom scope", "describe channeltypes --scope app1.app2 hello_world").
		WithExample("Display the state of the given channelType by the path", "describe channeltypes app1.app2.hello_world").
		WithAliases([]string{"ct"}).
		WithCommonFlags().
		ExactArgs(1, displayChannelTypeState)

	describeNode := cmd.NewCmd("nodes").
		WithDescription("retrieves the full state of the node from a given namespace").
		WithExample("Display the state of the given node on the default scope", "describe node hello_world").
		WithExample("Display the state of the given node on a custom scope", "describe node --scope app1.app2 hello_world").
		WithExample("Display the state of the given node by the path", "describe node app1.app2.hello_world").
		WithAliases([]string{"n"}).
		ExactArgs(1, displayNodeState)

	describeCmd.AddCommand(describeApp)
	describeCmd.AddCommand(describeChannel)
	describeCmd.AddCommand(describeChannelType)
	describeCmd.AddCommand(describeNode)

	return describeCmd
}

func displayAppState(_ context.Context, out io.Writer, args []string) error {
	client := getClient()

	scope, err := getScope()
	if err != nil {
		return err
	}

	if !utils.IsValidScope(args[0]) {
		fmt.Println("invalid args")
		return ierrors.NewError().Build()
	}

	separator := ""
	if scope != "" {
		separator = "."
	}
	path := scope + separator + args[0]

	app, err := client.Apps().Get(context.Background(), path)
	if err != nil {
		fmt.Println(err)
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

	path := scope
	var chName string

	if err := utils.StructureNameIsValid(args[0]); err != nil {
		if !utils.IsValidScope(args[0]) {
			return ierrors.NewError().Build()
		}

		newScope, lastName, err := utils.RemoveLastPartInScope(args[0])
		if err != nil {
			return ierrors.NewError().Build()
		}

		separator := ""
		if scope != "" {
			separator = "."
		}

		path = path + separator + newScope
		chName = lastName

	} else {
		chName = args[0]
	}

	channel, err := client.Channels().Get(context.Background(), path, chName)
	if err != nil {
		fmt.Println(err)
		return ierrors.NewError().Build()
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

	path := scope
	var ctName string

	if err := utils.StructureNameIsValid(args[0]); err != nil {
		if !utils.IsValidScope(args[0]) {
			return ierrors.NewError().Build()
		}

		newScope, lastName, err := utils.RemoveLastPartInScope(args[0])
		if err != nil {
			return ierrors.NewError().Build()
		}

		separator := ""
		if scope != "" {
			separator = "."
		}

		path = path + separator + newScope
		ctName = lastName

	} else {
		ctName = args[0]
	}

	channelType, err := client.ChannelTypes().Get(context.Background(), path, ctName)
	if err != nil {
		fmt.Println(err)
		return ierrors.NewError().Build()
	}
	utils.PrintChannelTypeTree(channelType)

	return nil
}

func displayNodeState(_ context.Context, out io.Writer, args []string) error {
	return nil
}

func getClient() *client.Client {
	url := "http://127.0.0.1:8080" // Here it will get from viper

	client := client.Client{
		HTTPClient: request.NewClient().BaseURL(url).Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build(),
	}
	return &client
}

func getScope() (string, error) {
	defaultScope := "" // Here it will get from viper
	scope := defaultScope

	if cmd.InsprOptions.Scope != "" {
		if utils.IsValidScope(cmd.InsprOptions.Scope) {
			scope = cmd.InsprOptions.Scope
		} else {
			fmt.Println("invalid scope")
			return "", ierrors.NewError().Build()
		}
	}

	return scope, nil
}
