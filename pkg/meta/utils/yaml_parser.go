package utils

import (
	"github.com/inspr/inspr/cmd/sidecars"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"gopkg.in/yaml.v2"
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

// YamlToApp - deserializes the yaml to a meta.App struct
func YamlToApp(bytes []byte) (*meta.App, error) {
	app := &meta.App{
		Meta: meta.Metadata{Annotations: make(map[string]string)},
	}

	if err := yaml.Unmarshal(bytes, &app); err != nil {
		return app, ierrors.NewError().Message("Error parsing the file").Build()
	}

	if app.Meta.Name == "" {
		return &meta.App{}, ierrors.NewError().Message("dapp without name").Build()
	}

	return app, nil
}

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
				Message("Error parsing the file").
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

// YamlToKafkaConfig - deserializes the yaml to a sidecars.KafkaConfig struct
func YamlToKafkaConfig(bytes []byte) (sidecars.KafkaConfig, error) {
	var config sidecars.KafkaConfig

	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return sidecars.KafkaConfig{},
			ierrors.
				NewError().
				InvalidType().
				Message("Error parsing the file").
				Build()
	}

	return config, nil
}
