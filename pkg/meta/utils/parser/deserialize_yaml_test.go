package utils

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	yaml "gopkg.in/yaml.v2"
)

const (
	fileName = "channel_yaml_test.yaml"
)

func TestYamlToChannel(t *testing.T) {
	yamlString, mockChannel := createYaml()

	// creates a file with the expected syntax
	ioutil.WriteFile(
		fileName,
		[]byte(yamlString),
		os.ModePerm,
	)

	// reads file created
	f, _ := os.Open(fileName)
	channel, err := YamlToChannel(f)
	if err != nil {
		t.Errorf("unexpected error -> got %v, expected %v", err, nil)
	}

	// uses cmp Equal to not evaluate comparison between maps
	if cmp.Equal(
		channel,
		mockChannel,
		cmp.Options{
			cmp.FilterValues(func(x, y interface{}) bool {
				vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
				flag := (vx.IsValid() && vy.IsValid() && vx.Type() == vy.Type()) && (vx.Kind() == reflect.Map)
				return flag
			}, cmp.Comparer(func(_, _ interface{}) bool { return true })),

			cmp.FilterValues(func(x, y interface{}) bool { return true },
				cmp.Comparer(func(x, y interface{}) bool {
					return reflect.DeepEqual(x, y)
				}),
			),
		}) {
		t.Errorf("unexpected error -> got %v, expected %v", channel, mockChannel)
	}
	os.Remove(fileName)
}

func TestNonExistantFile(t *testing.T) {
	// reads file created
	f, _ := os.Open(fileName)
	_, err := YamlToChannel(f)
	if err == nil {
		t.Errorf("unexpected error -> got %v, expected %v", err, "error")
	}
}

// CreateYaml - creates an yaml example
func createYaml() (string, meta.Channel) {
	channel := meta.Channel{
		Meta: meta.Metadata{
			Name:        "mock_name",
			Reference:   "mock_reference",
			Annotations: make(map[string]string),
			Parent:      "mock_parent",
			SHA256:      "mock_sha256",
		},
		Spec:          meta.ChannelSpec{Type: "mock_type"},
		ConnectedApps: []string{"a", "b", "c"},
	}
	data, _ := yaml.Marshal(&channel)
	return string(data), channel
}
