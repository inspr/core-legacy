package cli

import (
	"reflect"
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
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "basic_read_test",
			wantErr: false,
		},
		{
			name:    "want_error",
			wantErr: true,
		},
	}

	initViperConfig()   // inits viper
	createViperConfig() // creates the config in the system in case it doesn't exists
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				viper.SetConfigFile("/etc/")
			}
			if err := readViperConfig(); (err != nil) != tt.wantErr {
				t.Errorf("readViperConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got := viper.Get(configCurrentScope); !tt.wantErr && !reflect.DeepEqual(got, "") {
				t.Errorf("readViperConfig() -> want = %v, got %v", "", got)
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
				key:   configCurrentScope,
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

	initViperConfig() // inits viper
	readViperConfig() // reads the current values of the viper config

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				viper.SetConfigFile("/etc/")
			}

			if err := changeViperValues(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("changeViperValues() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got := viper.Get(tt.args.key); !reflect.DeepEqual(got, tt.args.value) && !tt.wantErr {
				t.Errorf("viper.Get(key) got = %v, want %v", got, tt.args.value)
			}
		})
	}
}
