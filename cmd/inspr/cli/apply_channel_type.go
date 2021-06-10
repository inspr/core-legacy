package cli

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/inspr/inspr/pkg/cmd"
	cliutils "github.com/inspr/inspr/pkg/cmd/utils"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
	utils "github.com/inspr/inspr/pkg/meta/utils/parser"
)

// NewApplyType receives a controller TypeInterface and calls it's methods
// depending on the flags values
func NewApplyType() RunMethod {
	return func(data []byte, out io.Writer) error {
		c := cliutils.GetCliClient().Types()
		// unmarshal into a Type
		insprType, err := utils.YamlToType(data)
		if err != nil {
			return err
		}

		if schemaNeedsInjection(insprType.Schema) {
			insprType.Schema, err = injectedSchema(insprType.Schema)
		}
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
