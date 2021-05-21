package parser

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/inspr/inspr/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

const fileNameChannel = "channel_yaml_test.yaml"

func TestYamlToChannel(t *testing.T) {

	yamlString, mockChannel := createYaml()
	// creates a file with the expected syntax
	ioutil.WriteFile(
		fileNameChannel,
		[]byte(yamlString),
		os.ModePerm,
	)

	// reads file created
	bytes, err := ioutil.ReadFile(fileNameChannel)
	if err != nil {
		t.Errorf("couldn't read file")
	}

	channel, err := YamlToChannel(bytes)
	if err != nil {
		t.Errorf("YamlToChannel() error -> got %v, expected %v", err, nil)
	}

	// uses cmp Equal to not evaluate comparison between maps
	if !cmp.Equal(
		channel,
		mockChannel,
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
		t.Errorf("unexpected error -> got %v, expected %v", channel, mockChannel)
	}
	os.Remove(fileNameChannel)
}

func TestIncorrectChannelYaml(t *testing.T) {
	t.Run("", func(t *testing.T) {
		bytes, _ := ioutil.ReadFile("mock_incorrect.yaml")

		_, err := YamlToChannel(bytes)
		if err == nil {
			t.Errorf("expected %v, received %v\n", errors.New("channel without name").Error(), err)
		}
	})
}

func TestNonExistentChannelFile(t *testing.T) {
	// reads file created
	bytes := []byte{1}
	_, err := YamlToChannel(bytes)
	if err == nil {
		t.Errorf("expected -> %v, expected %v", err, "error")
	}
}

// CreateYaml - creates an yaml example
func createYaml() (string, meta.Channel) {
	channel := meta.Channel{
		Meta: meta.Metadata{
			Name:        "mock_name",
			Reference:   "mock_reference",
			Annotations: map[string]string{},
			Parent:      "mock_parent",
			UUID:        "mock_sha256",
		},
		Spec:          meta.ChannelSpec{Type: "mock_type"},
		ConnectedApps: []string{"a", "b", "c"},
	}
	data, _ := yaml.Marshal(&channel)
	return string(data), channel
}
