package utils

import (
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

// YamlToType - deserializes the yaml to a meta.Type struct
func YamlToType(bytes []byte) (meta.Type, error) {
	insprType := meta.Type{
		Meta: meta.Metadata{Annotations: make(map[string]string)},
	}

	if err := yaml.Unmarshal(bytes, &insprType); err != nil {
		return meta.Type{},
			ierrors.
				NewError().
				InvalidType().
				Message("error parsing type yaml file").
				Build()
	}

	if insprType.Meta.Name == "" {
		return meta.Type{},
			ierrors.
				NewError().
				InvalidName().
				Message("type without name").
				Build()
	}

	return insprType, nil
}
