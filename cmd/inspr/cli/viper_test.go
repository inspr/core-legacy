package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
)

func Test_initConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "sets_default_values",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initConfig()
			scope := viper.Get(configCurrentScope)
			if scope != "./app1/app2" {
				t.Errorf("viper's default scope, expected %v, got %v", "./app1/app2", scope)
			}
		})
	}
}

func Test_createConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "check_and_create_folder",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createConfig(); (err != nil) != tt.wantErr {
				t.Errorf("createConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_readConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantOut string
		wantErr bool
	}{
		{
			name:    "basic_read",
			wantOut: "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := readConfig(out); (err != nil) != tt.wantErr {
				t.Errorf("readConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("readConfig() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
