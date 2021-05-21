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

const fileNameAlias = "alias_test.yaml"

func TestYamlToAlias(t *testing.T) {
	yamlString, mockAlias := createAliasYaml()
	// creates a file with the expected syntax
	ioutil.WriteFile(
		fileNameAlias,
		[]byte(yamlString),
		os.ModePerm,
	)

	// reads file created
	bytes, err := ioutil.ReadFile(fileNameAlias)
	if err != nil {
		t.Errorf("couldn't read file")
	}

	alias, err := YamlToAlias(bytes)
	if err != nil {
		t.Errorf("YamlToAlias() error -> got %v, expected %v", err, nil)
	}

	// uses cmp Equal to not evaluate comparison between maps
	if !cmp.Equal(
		alias,
		mockAlias,
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
		t.Errorf("unexpected error -> got %v, expected %v", alias, mockAlias)
	}
	os.Remove(fileNameAlias)
}

func TestIncorrectAliasYaml(t *testing.T) {
	t.Run("", func(t *testing.T) {
		bytes, _ := ioutil.ReadFile("mock_incorrect.yaml")

		_, err := YamlToAlias(bytes)
		if err == nil {
			t.Errorf("expected %v, received %v\n", errors.New("alias without name").Error(), err)
		}
	})
}

func TestNonExistantAliasFile(t *testing.T) {
	// reads file created
	bytes := []byte{1}
	_, err := YamlToAlias(bytes)
	if err == nil {
		t.Errorf("expected -> %v, expected %v", err, "error")
	}
}

// CreateYaml - creates an yaml example
func createAliasYaml() (string, *meta.Alias) {
	alias := &meta.Alias{
		Meta: meta.Metadata{
			Name: "mock_name",
		},
		Target: "mock_target",
	}
	data, _ := yaml.Marshal(&alias)
	return string(data), alias
}
