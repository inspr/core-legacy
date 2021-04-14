package cli

import (
	"context"
	"io"

	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	metautils "gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	utils "gitlab.inspr.dev/inspr/core/pkg/meta/utils/parser"
)

// NewApplyAlias receives a controller AliasInterface and calls it's methods
// depending on the flags values
func NewApplyAlias() RunMethod {
	return func(data []byte, out io.Writer) error {
		c := cliutils.GetCliClient().Alias()
		// unmarshal into a channel
		alias, err := utils.YamlToAlias(data)
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

		parentPath, err := metautils.JoinScopes(scope, alias.Meta.Parent)
		if err != nil {
			return err
		}

		// creates or updates it
		if flagIsUpdate {
			log, err = c.Update(context.Background(), parentPath, alias.Meta.Name, &alias, flagDryRun)
		} else {
			log, err = c.Create(context.Background(), parentPath, alias.Meta.Name, &alias, flagDryRun)
		}

		if err != nil {
			return err
		}

		// prints differences
		log.Print(out)

		return nil
	}
}
