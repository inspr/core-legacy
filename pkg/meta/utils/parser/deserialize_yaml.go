package utils

import (
	"errors"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

// YamlToChannel - deserializes the yaml to a meta.Channel struct
func YamlToChannel(bytes []byte) (meta.Channel, error) {
	channel := meta.Channel{}

	if err := yaml.Unmarshal(bytes, &channel); err != nil {
		return channel, errors.New("Error parsing the file")
	}

	if channel.Meta.Name == "" {
		return meta.Channel{}, errors.New("channel without name")
	}

	return channel, nil
}
