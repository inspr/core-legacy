package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func Test_initViperConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "sets_default_values",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitViperConfig()
			scope := viper.Get(configScope)
			if scope != defaultValues[configScope] {
				t.Errorf("viper's scope, expected %v, got %v",
					scope,
					defaultValues[configScope])
			}

			ip := viper.Get(configServerIP)
			if ip != defaultValues[configServerIP] {

				t.Errorf("viper's serverip, expected %v, got %v",
					ip,
					defaultValues[configServerIP])
			}
		})
	}
}

func Test_changeViperValues(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "changing_scope",
			args: args{
				key:   configScope,
				value: "new_scope",
			},
			wantErr: false,
		},
		{
			name: "changing_IP",
			args: args{
				key:   configServerIP,
				value: "XXX.YYY.ZZZ.0",
			},
			wantErr: false,
		},
		{
			name: "error_writing",
			args: args{
				key:   configServerIP,
				value: nil,
			},
			wantErr: true,
		},
	}

	// read mock config values
	folder := setupViperTest(t)
	os.Setenv("HOME", folder)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				os.Setenv("HOME", "/etc/")
			}

			// reads the current values of the viper config
			ReadDefaultConfig()

			err := ChangeViperValues(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("changeViperValues() error = %v, wantErr %v",
					err,
					tt.wantErr)
			}

			got := viper.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.args.value) && !tt.wantErr {
				t.Errorf("viper.Get(key) got = %v, want %v",
					got,
					tt.args.value)
			}
		})
	}
}

func Test_existingKeys(t *testing.T) {
	InitViperConfig()
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "all_keys",
			want: []string{"scope", "serverip"},
		},
	}

	// read mock config values

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExistingKeys()

			// doesn't return in a specific order, so we use maps to compare
			receivedValues := make(map[string]bool)
			for _, k := range got {
				receivedValues[k] = true
			}

			// checking values
			for _, k := range tt.want {
				if receivedValues[k] == false {
					t.Errorf("existingKeys() => %v doesn't exist but is expected",
						k)
				}
			}
		})
	}
}

/// test utils functions

func setupViperTest(t *testing.T) string {
	folder := t.TempDir()
	config := struct {
		ServerIP string `yaml:"serverip"`
		Scope    string `yaml:"scope"`
	}{
		ServerIP: "http://localhost:8080",
	}
	mars, _ := yaml.Marshal(config)
	err := os.MkdirAll(filepath.Join(folder, ".inspr"), 0755)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = os.WriteFile(filepath.Join(folder, ".inspr", "config"), mars, 0644)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return folder
}
