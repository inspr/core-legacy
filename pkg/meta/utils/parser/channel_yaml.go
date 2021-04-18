package utils

import (
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

// YamlToChannel - deserializes the yaml to a meta.Channel struct
func YamlToChannel(bytes []byte) (meta.Channel, error) {
	channel := meta.Channel{
		Meta: meta.Metadata{Annotations: make(map[string]string)},
	}

	if err := yaml.Unmarshal(bytes, &channel); err != nil {
		return channel, ierrors.NewError().Message("Error parsing the file").Build()
	}

	if channel.Meta.Name == "" {
		return meta.Channel{}, ierrors.NewError().Message("channel without name").Build()
	}

	return channel, nil
}
