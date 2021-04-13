package utils

import (
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

// YamlToAlias - deserializes the yaml to a meta.Alias struct
func YamlToAlias(bytes []byte) (meta.Alias, error) {
	var alias *meta.Alias

	if err := yaml.Unmarshal(bytes, &alias); err != nil {
		return *alias, ierrors.NewError().Message("Error parsing the file").Build()
	}

	if alias.Meta.Name == "" {
		return *alias, ierrors.NewError().Message("alias without name").Build()
	}

	return *alias, nil
}
