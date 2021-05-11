package utils

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

const (
	typefile = "type_yaml_test.yaml"
)

func TestYamlToType(t *testing.T) {

	yamlString, mockCT := createTypeYaml()
	// creates a file with the expected syntax
	ioutil.WriteFile(
		typefile,
		[]byte(yamlString),
		os.ModePerm,
	)

	// reads file created
	bytes, err := ioutil.ReadFile(typefile)
	if err != nil {
		t.Errorf("couldn't read file")
	}

	channel, err := YamlToType(bytes)
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
	os.Remove(typefile)
}

func TestIncorrecttypeYaml(t *testing.T) {
	t.Run("", func(t *testing.T) {
		bytes, _ := ioutil.ReadFile("mock_incorrect.yaml")

		_, err := YamlToType(bytes)
		if err == nil {
			t.Errorf("expected %v, received %v\n", errors.New("channel without name").Error(), err)
		}
	})
}

func TestNonExistenttypefile(t *testing.T) {
	// reads file created
	bytes := []byte{1}
	_, err := YamlToType(bytes)
	if err == nil {
		t.Errorf("expected -> %v, expected %v", err, "error")
	}
}

// CreateYaml - creates an yaml example
func createTypeYaml() (string, meta.Type) {
	ct := meta.Type{
		Meta: meta.Metadata{
			Name:        "mock_name",
			Reference:   "mock_reference",
			Annotations: map[string]string{},
			Parent:      "mock_parent",
			UUID:        "mock_sha256",
		},
		Schema:            "mock_schema",
		ConnectedChannels: []string{"mock_chan1", "mock_chan2"},
	}
	data, _ := yaml.Marshal(&ct)
	return string(data), ct
}
