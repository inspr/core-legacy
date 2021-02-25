package cli

import (
	"context"
	"io"

	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/controller"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	metautils "gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	utils "gitlab.inspr.dev/inspr/core/pkg/meta/utils/parser"
)

// NewApplyApp receives a controller AppInterface and calls it's methods
// depending on the flags values
func NewApplyApp(c controller.AppInterface) RunMethod {
	return func(data []byte, out io.Writer) error {
		// unmarshal into an app
		app, err := utils.YamlToApp(data)
		if err != nil {
			return err
		}

		if len(app.Spec.ChannelTypes) > 0 {
			err = schemaInjection(app.Spec.ChannelTypes)
			if err != nil {
				return err
			}
		}

		if len(app.Spec.Apps) > 0 {
			err = recursiveSchemaInjection(app.Spec.Apps)
			if err != nil {
				return err
			}
		}

		flagDryRun := cmd.InsprOptions.DryRun
		flagIsUpdate := cmd.InsprOptions.Update

		var log diff.Changelog
		query, err := metautils.JoinScopes(cmd.InsprOptions.Scope, app.Meta.Parent)
		if err != nil {
			return err
		}
		// creates or updates it
		if flagIsUpdate {
			log, err = c.Update(context.Background(), query, &app, flagDryRun)
		} else {
			log, err = c.Create(context.Background(), query, &app, flagDryRun)
		}

		if err != nil {
			return err
		}

		// prints differences
		log.Print(out)

		return nil
	}
}

func schemaInjection(ctypes map[string]*meta.ChannelType) error {
	var err error
	for _, ctype := range ctypes {
		if schemaNeedsInjection(ctype.Schema) {
			ctype.Schema, err = injectedSchema(ctype.Schema)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func recursiveSchemaInjection(apps map[string]*meta.App) error {
	var err error
	for _, app := range apps {
		if len(app.Spec.ChannelTypes) > 0 {
			err = schemaInjection(app.Spec.ChannelTypes)
			if err != nil {
				return err
			}
		}
		if len(app.Spec.Apps) > 0 {
			err = recursiveSchemaInjection(app.Spec.Apps)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
