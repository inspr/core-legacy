package cli

import (
	"testing"

	"github.com/spf13/viper"
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
			initViperConfig()
			scope := viper.Get(configCurrentScope)
			if scope != "" {
				t.Errorf("viper's default scope, expected %v, got %v", "", scope)
			}
		})
	}
}

func Test_createViperConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "check_and_create_folder",
			wantErr: false,
		},
	}
	initViperConfig()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createViperConfig(); (err != nil) != tt.wantErr {
				t.Errorf("createConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_readViperConfig(t *testing.T) {
	name := "basic_read_test"
	wantErr := false
	initViperConfig()   // inits viper
	createViperConfig() // creates the config in the system in case it doesn't exists

	// tests
	t.Run(name, func(t *testing.T) {

		if err := readViperConfig(); (err != nil) != wantErr {
			t.Errorf("readConfig() error = %v, wantErr %v", err, wantErr)
		}
		scope := viper.Get(configCurrentScope)
		if scope != defaultValues[configCurrentScope] {
			t.Errorf("readConfig() -> scope, error = %v, wantErr %v", scope, defaultValues[configCurrentScope])
		}

		ip := viper.Get(configServerIP)
		if ip != defaultValues[configServerIP] {
			t.Errorf("readConfig() -> scope, error = %v, wantErr %v", scope, defaultValues[configServerIP])
		}
	})

}
