package cli

import (
	"context"
	"io"

	"github.com/inspr/inspr/pkg/cmd"
	cliutils "github.com/inspr/inspr/pkg/cmd/utils"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
	"gopkg.in/yaml.v2"
)

// NewApplyChannel receives a controller ChannelInterface and calls it's methods
// depending on the flags values
func NewApplyChannel() RunMethod {
	return func(data []byte, out io.Writer) error {
		c := cliutils.GetCliClient().Channels()
		var channel meta.Channel = meta.Channel{
			Meta: meta.Metadata{Annotations: make(map[string]string)},
		}

		// unmarshal into a channel
		if err := yaml.Unmarshal(data, &channel); err != nil {
			return err
		}
		if channel.Meta.Name == "" {
			return ierrors.NewError().Message("channel without name").Build()
		}

		flagDryRun := cmd.InsprOptions.DryRun
		flagIsUpdate := cmd.InsprOptions.Update

		var log diff.Changelog

		scope, err := cliutils.GetScope()
		if err != nil {
			return err
		}

		parentScope, err := metautils.JoinScopes(scope, channel.Meta.Parent)
		if err != nil {
			return err
		}

		// creates or updates it
		if flagIsUpdate {
			log, err = c.Update(context.Background(), parentScope, &channel, flagDryRun)
		} else {
			log, err = c.Create(context.Background(), parentScope, &channel, flagDryRun)
		}

		if err != nil {
			return err
		}

		// prints differences
		log.Print(out)

		return nil
	}
}
