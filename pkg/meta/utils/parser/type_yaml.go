package utils

import (
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

// YamlToType - deserializes the yaml to a meta.Type struct
func YamlToType(bytes []byte) (meta.Type, error) {
	Type := meta.Type{
		Meta: meta.Metadata{Annotations: make(map[string]string)},
	}

	if err := yaml.Unmarshal(bytes, &Type); err != nil {
		return meta.Type{},
			ierrors.
				NewError().
				InvalidType().
				Message("Error parsing the file").
				Build()
	}

	if Type.Meta.Name == "" {
		return meta.Type{},
			ierrors.
				NewError().
				InvalidName().
				Message("Type without name").
				Build()
	}

	return Type, nil
}
