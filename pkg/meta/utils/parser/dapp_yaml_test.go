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

const appFile = "dapp_yaml_test.yaml"

func TestYamlToApp(t *testing.T) {

	yamlString, mockApp := createDAppYaml()
	// creates a file with the expected syntax
	ioutil.WriteFile(
		appFile,
		[]byte(yamlString),
		os.ModePerm,
	)

	// reads file created
	bytes, err := ioutil.ReadFile(appFile)
	if err != nil {
		t.Errorf("couldn't read file")
	}

	app, err := YamlToApp(bytes)
	if err != nil {
		t.Errorf("YamlToApp() error -> got %v, expected %v", err, nil)
	}

	// uses cmp Equal to not evaluate comparison between maps
	if !cmp.Equal(
		*app,
		mockApp,
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
		t.Errorf("unexpected error -> got %v, expected %v", app, mockApp)
	}
	os.Remove(appFile)
}

func TestIncorrectDAppYaml(t *testing.T) {
	t.Run("", func(t *testing.T) {
		bytes, _ := ioutil.ReadFile("mock_incorrect.yaml")

		_, err := YamlToApp(bytes)
		if err == nil {
			t.Errorf("expected %v, received %v\n", errors.New("dapp without name").Error(), err)
		}
	})
}

func TestNonExistentDAppFile(t *testing.T) {
	// reads file created
	bytes := []byte{1}
	_, err := YamlToApp(bytes)
	if err == nil {
		t.Errorf("expected -> %v, expected %v", err, "error")
	}
}

// CreateYaml - creates an yaml example
func createDAppYaml() (string, meta.App) {
	app := meta.App{
		Meta: meta.Metadata{
			Name:        "mock_name",
			Reference:   "mock_reference",
			Annotations: map[string]string{},
			Parent:      "mock_parent",
			UUID:        "mock_sha256",
		},
		Spec: meta.AppSpec{},
	}
	data, _ := yaml.Marshal(&app)
	return string(data), app
}
