package cli

import (
	"context"
	"io"

	"gopkg.in/yaml.v2"
	"inspr.dev/inspr/pkg/cmd"
	cliutils "inspr.dev/inspr/pkg/cmd/utils"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	metautils "inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

// NewApplyApp receives a controller AppInterface and calls it's methods
// depending on the flags values
func NewApplyApp() RunMethod {
	return func(data []byte, out io.Writer) error {
		c := cliutils.GetCliClient().Apps()
		var app meta.App = meta.App{
			Meta: meta.Metadata{Annotations: make(map[string]string)},
		}

		// unmarshal into an app

		if err := yaml.Unmarshal(data, &app); err != nil {
			return err
		}
		if app.Meta.Name == "" {
			return ierrors.New("dapp without name")
		}

		err := recursiveSchemaInjection(&app)
		if err != nil {
			return err
		}

		flagDryRun := cmd.InsprOptions.DryRun
		flagIsUpdate := cmd.InsprOptions.Update

		var log diff.Changelog
		scope, err := metautils.JoinScopes(cmd.InsprOptions.Scope, app.Meta.Parent)
		if err != nil {
			return err
		}
		// creates or updates it
		if flagIsUpdate {
			updateQuery, errQuery := metautils.JoinScopes(scope, app.Meta.Name)
			if errQuery != nil {
				return errQuery
			}
			log, err = c.Update(context.Background(), updateQuery, &app, flagDryRun)
		} else {
			log, err = c.Create(context.Background(), scope, &app, flagDryRun)
		}

		if err != nil {
			return err
		}

		// prints differences
		log.Print(out)

		return nil
	}
}

func schemaInjection(types map[string]*meta.Type) error {
	var err error
	for typeName, insprType := range types {
		insprType.Meta.Name = typeName
		if schemaNeedsInjection(insprType.Schema) {
			insprType.Schema, err = injectedSchema(insprType.Schema)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func recursiveSchemaInjection(app *meta.App) error {
	var err error

	if len(app.Spec.Types) > 0 {
		err = schemaInjection(app.Spec.Types)
		if err != nil {
			return err
		}
	}

	for chName, channel := range app.Spec.Channels {
		channel.Meta.Name = chName
	}

	for aliasName, alias := range app.Spec.Aliases {
		alias.Meta.Name = aliasName
	}

	if app.Spec.Node.Spec.Image != "" {
		if app.Spec.Node.Spec.Replicas == 0 {
			app.Spec.Node.Spec.Replicas = 1
		}
	}

	for appName, childApp := range app.Spec.Apps {
		childApp.Meta.Name = appName
		err = recursiveSchemaInjection(childApp)
		if err != nil {
			return err
		}
	}

	return nil
}
