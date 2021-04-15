package utils

import (
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

// YamlToApp - deserializes the yaml to a meta.App struct
func YamlToApp(bytes []byte) (meta.App, error) {
	app := meta.App{
		Meta: meta.Metadata{Annotations: make(map[string]string)},
	}

	if err := yaml.Unmarshal(bytes, &app); err != nil {
		return app, ierrors.NewError().Message("Error parsing the file").Build()
	}

	if app.Meta.Name == "" {
		return meta.App{}, ierrors.NewError().Message("dapp without name").Build()
	}

	return app, nil
}
