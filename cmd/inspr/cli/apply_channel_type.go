package cli

import (
	"context"
	"io"

	"gitlab.inspr.dev/inspr/core/pkg/controller"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	utils "gitlab.inspr.dev/inspr/core/pkg/meta/utils/parser"
)

// ApplyChannelType is of the type RunMethod, it calls the pkg/controller/client functions.
type ApplyChannelType RunMethod

// NewApplyChannelType receives a controller ChannelTypeInterface and calls it's methods
// depending on the flags values
func NewApplyChannelType(c controller.ChannelTypeInterface) RunMethod {
	return func(data []byte, out io.Writer) error {
		// unmarshal into a channel
		channel, err := utils.YamlToChannelType(data)
		if err != nil {
			return err
		}

		flagDryRun := false
		flagIsUpdate := false

		var log diff.Changelog
		// creates or updates it
		if flagIsUpdate {
			log, err = c.Update(context.Background(), channel.Meta.Parent, &channel, flagDryRun)
		} else {
			log, err = c.Create(context.Background(), channel.Meta.Parent, &channel, flagDryRun)
		}

		if err != nil {
			return err
		}

		// prints differences
		log.Print(out)

		return nil
	}
}
