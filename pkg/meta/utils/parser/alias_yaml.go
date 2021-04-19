package utils

import (
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

// YamlToAlias - deserializes the yaml to a meta.Alias struct
func YamlToAlias(bytes []byte) (*meta.Alias, error) {
	var alias *meta.Alias

	if err := yaml.Unmarshal(bytes, &alias); err != nil {
		return nil, ierrors.NewError().Message(err.Error()).Build()
	}

	if alias.Meta.Name == "" {
		return nil, ierrors.NewError().Message("alias without name").Build()
	}

	return alias, nil
}
