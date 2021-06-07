package cli

import (
	"context"
	"io"

	"github.com/inspr/inspr/pkg/cmd"
	cliutils "github.com/inspr/inspr/pkg/cmd/utils"
	"github.com/inspr/inspr/pkg/meta"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
	"gopkg.in/yaml.v2"
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
