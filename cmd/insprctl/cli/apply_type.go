package cli

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
	"inspr.dev/inspr/pkg/cmd"
	cliutils "inspr.dev/inspr/pkg/cmd/utils"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	metautils "inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

// NewApplyType receives a controller TypeInterface and calls it's methods
// depending on the flags values
func NewApplyType() RunMethod {
	return func(data []byte, out io.Writer) error {
		c := cliutils.GetCliClient().Types()
		var insprType meta.Type = meta.Type{
			Meta: meta.Metadata{Annotations: make(map[string]string)},
		}

		// unmarshal into a Type
		if err := yaml.Unmarshal(data, &insprType); err != nil {
			return err
		}
		if insprType.Meta.Name == "" {
			return ierrors.New("type without name")
		}

		if schemaNeedsInjection(insprType.Schema) {
			var err error
			insprType.Schema, err = injectedSchema(insprType.Schema)
			if err != nil {
				return err
			}
		} else if !IsJSON(insprType.Schema) {
			return ierrors.New("invalid type schema")
		}

		flagDryRun := cmd.InsprOptions.DryRun
		flagIsUpdate := cmd.InsprOptions.Update

		var log diff.Changelog

		scope, err := cliutils.GetScope()
		if err != nil {
			return err
		}

		parentScope, err := metautils.JoinScopes(scope, insprType.Meta.Parent)
		if err != nil {
			return err
		}

		// creates or updates it
		if flagIsUpdate {
			log, err = c.Update(context.Background(), parentScope, &insprType, flagDryRun)
		} else {
			log, err = c.Create(context.Background(), parentScope, &insprType, flagDryRun)
		}

		if err != nil {
			return err
		}

		// prints differences
		log.Print(out)

		return nil
	}
}

func schemaNeedsInjection(schema string) bool {
	_, err := os.Stat(schema)
	if !os.IsNotExist(err) &&
		(filepath.Ext(schema) == ".schema" || filepath.Ext(schema) == ".avsc") {
		// file exists and has the right extension
		return true
	}
	return false
}

func injectedSchema(path string) (string, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	schema := string(file)

	return schema, nil
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
