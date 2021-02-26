package cli

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/controller"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	utils "gitlab.inspr.dev/inspr/core/pkg/meta/utils/parser"
)

// NewApplyChannelType receives a controller ChannelTypeInterface and calls it's methods
// depending on the flags values
func NewApplyChannelType(c controller.ChannelTypeInterface) RunMethod {
	return func(data []byte, out io.Writer) error {
		// unmarshal into a channelType
		channelType, err := utils.YamlToChannelType(data)
		if err != nil {
			return err
		}

		if schemaNeedsInjection(channelType.Schema) {
			channelType.Schema, err = injectedSchema(channelType.Schema)
		}

		flagDryRun := cmd.InsprOptions.DryRun
		flagIsUpdate := cmd.InsprOptions.Update

		var log diff.Changelog
		// creates or updates it
		if flagIsUpdate {
			log, err = c.Update(context.Background(), channelType.Meta.Parent, &channelType, flagDryRun)
		} else {
			log, err = c.Create(context.Background(), channelType.Meta.Parent, &channelType, flagDryRun)
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
		// file exists and has the right extention
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
