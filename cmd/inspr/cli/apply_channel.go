package cli

import (
	"context"
	"io"

	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/controller"
	metautils "gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	utils "gitlab.inspr.dev/inspr/core/pkg/meta/utils/parser"
)

// NewApplyChannel receives a controller ChannelInterface and calls it's methods
// depending on the flags values
func NewApplyChannel(c controller.ChannelInterface) RunMethod {
	return func(data []byte, out io.Writer) error {
		// unmarshal into a channel
		channel, err := utils.YamlToChannel(data)
		if err != nil {
			return err
		}

		flagDryRun := cmd.InsprOptions.DryRun
		flagIsUpdate := cmd.InsprOptions.Update

		var log diff.Changelog

		scope, err := cliutils.GetScope()
		if err != nil {
			return err
		}

		parentPath, err := metautils.JoinScopes(scope, channel.Meta.Parent)
		if err != nil {
			return err
		}

		// creates or updates it
		if flagIsUpdate {
			log, err = c.Update(context.Background(), parentPath, &channel, flagDryRun)
		} else {
			log, err = c.Create(context.Background(), parentPath, &channel, flagDryRun)
		}

		if err != nil {
			return err
		}

		// prints differences
		log.Print(out)

		return nil
	}
}
