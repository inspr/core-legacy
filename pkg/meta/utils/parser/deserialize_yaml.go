package utils

import (
	"errors"
	"io/ioutil"
	"os"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

// YamlToChannel - deserializes the yaml to a meta.Channel struct
func YamlToChannel(f *os.File) (meta.Channel, error) {
	channel := meta.Channel{}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return channel, errors.New("Error reading the file")
	}
	err = yaml.Unmarshal(bytes, &channel)
	if err != nil {
		return channel, errors.New("Error parsing the file")
	}

	return channel, nil
}
