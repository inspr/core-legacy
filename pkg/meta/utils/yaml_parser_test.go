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
	fileNameAlias   = "alias_test.yaml"
	fileNameChannel = "channel_yaml_test.yaml"
	fileNameDapp    = "dapp_yaml_test.yaml"
	fileNameType    = "type_yaml_test.yaml"

	// used for creating invalid yaml files locally, that should be deleted after each test
	incorrectYamlContent = `
incorrect_tag_meta:
	name: mock_name
	reference: mock_reference
	annotations:
		mock_map: mock_value
	parent: mock_parent
	sha256: mock_sha256
incorrect_tag_spec:
	type: mock_type
incorrect_tag_connectedapps:
- a
- b
- c
`
)

// ---------------------- MOCKING section ----------------------

// createInvalidYaml creates an invalid yaml
func createInvalidYaml() {
	ioutil.WriteFile("mock_incorrect.yaml", []byte(incorrectYamlContent), 0777)
}

// deleteInvalidYaml deletes the invalid yaml used for testing
func deleteInvalidYaml() {
	os.Remove("mock_incorrect.yaml")
}

// createAliasYaml - creates an yaml example
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

// createChannelYaml - creates channel an yaml example
func createChannelYaml() (string, meta.Channel) {
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

// createDAppYaml - creates an yaml example
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

// createTypeYaml - creates an yaml example
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

// ---------------------- ALIAS section ----------------------
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
	createInvalidYaml()
	defer deleteInvalidYaml()
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

func TestYamlToChannel(t *testing.T) {

	yamlString, mockChannel := createChannelYaml()
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

// ---------------------- CHANNEL section ----------------------
func TestIncorrectChannelYaml(t *testing.T) {
	createInvalidYaml()
	defer deleteInvalidYaml()
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

// ---------------------- DAPP section ----------------------
func TestYamlToApp(t *testing.T) {

	yamlString, mockApp := createDAppYaml()
	// creates a file with the expected syntax
	ioutil.WriteFile(
		fileNameDapp,
		[]byte(yamlString),
		os.ModePerm,
	)

	// reads file created
	bytes, err := ioutil.ReadFile(fileNameDapp)
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
	os.Remove(fileNameDapp)
}

func TestIncorrectDAppYaml(t *testing.T) {
	createInvalidYaml()
	defer deleteInvalidYaml()
	t.Run("", func(t *testing.T) {
		bytes, _ := ioutil.ReadFile("mock_incorrect.yaml")

		_, err := YamlToApp(bytes)
		if err == nil {
			t.Errorf("expected %v, received %v\n", errors.New("dapp without name").Error(), err)
		}
	})
}

func TestNonExistentDfileNameDapp(t *testing.T) {
	// reads file created
	bytes := []byte{1}
	_, err := YamlToApp(bytes)
	if err == nil {
		t.Errorf("expected -> %v, expected %v", err, "error")
	}
}

// ---------------------- TYPE section ----------------------
func TestYamlToType(t *testing.T) {

	yamlString, mockCT := createTypeYaml()
	// creates a file with the expected syntax
	ioutil.WriteFile(
		fileNameType,
		[]byte(yamlString),
		os.ModePerm,
	)

	// reads file created
	bytes, err := ioutil.ReadFile(fileNameType)
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
	os.Remove(fileNameType)
}

func TestIncorrecttypeYaml(t *testing.T) {
	createInvalidYaml()
	defer deleteInvalidYaml()
	t.Run("", func(t *testing.T) {
		bytes, _ := ioutil.ReadFile("mock_incorrect.yaml")

		_, err := YamlToType(bytes)
		if err == nil {
			t.Errorf("expected %v, received %v\n", errors.New("channel without name").Error(), err)
		}
	})
}

func TestNonExistentfileNameType(t *testing.T) {
	// reads file created
	bytes := []byte{1}
	_, err := YamlToType(bytes)
	if err == nil {
		t.Errorf("expected -> %v, expected %v", err, "error")
	}
}
