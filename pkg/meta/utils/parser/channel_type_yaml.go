package utils

import (
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

// YamlToChannelType - deserializes the yaml to a meta.ChannelType struct
func YamlToChannelType(bytes []byte) (meta.ChannelType, error) {
	channelType := meta.ChannelType{
		Meta: meta.Metadata{Annotations: make(map[string]string)},
	}

	if err := yaml.Unmarshal(bytes, &channelType); err != nil {
		return meta.ChannelType{},
			ierrors.
				NewError().
				InvalidChannelType().
				Message("Error parsing the file").
				Build()
	}

	if channelType.Meta.Name == "" {
		return meta.ChannelType{},
			ierrors.
				NewError().
				InvalidName().
				Message("channelType without name").
				Build()
	}

	return channelType, nil
}
