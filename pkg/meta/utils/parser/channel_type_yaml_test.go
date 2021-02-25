package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

const (
	channelTypeFile = "channel_type_yaml_test.yaml"
)

func TestYamlToChannelType(t *testing.T) {

	yamlString, mockCT := createChannelTypeYaml()
	// creates a file with the expected syntax
	ioutil.WriteFile(
		channelTypeFile,
		[]byte(yamlString),
		os.ModePerm,
	)

	// reads file created
	bytes, err := ioutil.ReadFile(channelTypeFile)
	if err != nil {
		t.Errorf("couldn't read file")
	}

	channel, err := YamlToChannelType(bytes)
	if err != nil {
		t.Errorf("YamlToChannel() error -> got %v, expected %v", err, nil)
	}

	// uses cmp Equal to not evaluate comparison between maps
	if !cmp.Equal(
		channel,
		mockCT,
		cmp.Options{
			cmp.FilterValues(func(x, y interface{}) bool {
				vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
				flag := (vx.IsValid() && vy.IsValid() && vx.Type() == vy.Type()) && (vx.Kind() == reflect.Map)
				return flag
			}, cmp.Comparer(func(_, _ interface{}) bool { return true })),

			// everything besides maps
			cmp.FilterValues(func(x, y interface{}) bool { return true },
				cmp.Comparer(func(x, y interface{}) bool {
					return reflect.DeepEqual(x, y)
				}),
			),
		}) {
		t.Errorf("unexpected error -> got %v, expected %v", channel, mockCT)
	}
	os.Remove(channelTypeFile)
}

func TestIncorrectCTypeYaml(t *testing.T) {
	t.Run("", func(t *testing.T) {
		bytes, _ := ioutil.ReadFile("mock_incorrect.yaml")

		_, err := YamlToChannelType(bytes)
		if err == nil {
			t.Errorf("expected %v, received %v\n", errors.New("channel without name").Error(), err)
		}
	})
}

func TestNonExistantCTypeFile(t *testing.T) {
	// reads file created
	bytes := []byte{1}
	_, err := YamlToChannelType(bytes)
	if err == nil {
		t.Errorf("expected -> %v, expected %v", err, "error")
	}
}

// CreateYaml - creates an yaml example
func createChannelTypeYaml() (string, meta.ChannelType) {
	ct := meta.ChannelType{
		Meta: meta.Metadata{
			Name:        "mock_name",
			Reference:   "mock_reference",
			Annotations: map[string]string{},
			Parent:      "mock_parent",
			SHA256:      "mock_sha256",
		},
		Schema:            "mock_schema",
		ConnectedChannels: []string{"mock_chan1", "mock_chan2"},
	}
	data, _ := yaml.Marshal(&ct)
	return string(data), ct
}
