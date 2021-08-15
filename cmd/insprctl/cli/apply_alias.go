package cli

import (
	"context"
	"io"

	"gopkg.in/yaml.v2"
	"inspr.dev/inspr/pkg/cmd"
	cliutils "inspr.dev/inspr/pkg/cmd/utils"
	"inspr.dev/inspr/pkg/meta"
	metautils "inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

// NewApplyAlias receives a controller AliasInterface and calls it's methods
// depending on the flags values
func NewApplyAlias() RunMethod {
	return func(data []byte, out io.Writer) error {
		c := cliutils.GetCliClient().Alias()
		var alias meta.Alias

		// unmarshal into a channel
		err := yaml.Unmarshal(data, &alias)
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

		parentScope, err := metautils.JoinScopes(scope, alias.Meta.Parent)
		if err != nil {
			return err
		}

		// creates or updates it
		if flagIsUpdate {
			log, err = c.Update(
				context.Background(),
				parentScope,
				alias.Meta.Name,
				&alias,
				flagDryRun,
			)
		} else {
			log, err = c.Create(context.Background(), parentScope, alias.Meta.Name, &alias, flagDryRun)
		}

		if err != nil {
			return err
		}

		// prints differences
		log.Print(out)

		return nil
	}
}
