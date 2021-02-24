package cli

import (
	"context"
	"io"

	"gitlab.inspr.dev/inspr/core/pkg/controller"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	utils "gitlab.inspr.dev/inspr/core/pkg/meta/utils/parser"
)

// RunMethod defines the method that will run for the component
type RunMethod func(data []byte, out io.Writer) error

// ApplyChannel is of the type RunMethod, it calls the pkg/controller/client functions.
type ApplyChannel RunMethod

// NewApplyChannel receives a controller ChannelInterface and calls it's methods
// depending on the flags values
func NewApplyChannel(c controller.ChannelInterface) RunMethod {
	return func(data []byte, out io.Writer) error {
		// unmarshal into a channel
		channel, err := utils.YamlToChannel(data)
		if err != nil {
			return err
		}

		flagDryRun := false
		flagIsUpdate := false

		var clog diff.Changelog
		// creates or updates it
		if flagIsUpdate {
			clog, err = c.Update(context.Background(), channel.Meta.Parent, &channel, flagDryRun)
		} else {
			clog, err = c.Create(context.Background(), channel.Meta.Parent, &channel, flagDryRun)
		}

		if err != nil {
			return err
		}

		// prints differences
		clog.Print(out)

		return nil
	}
}
