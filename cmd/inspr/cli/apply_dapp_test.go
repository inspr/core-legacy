package cli

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/controller"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func TestNewApplyApp(t *testing.T) {
	type args struct {
		c controller.AppInterface
	}
	tests := []struct {
		name string
		args args
		want RunMethod
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewApplyApp(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApplyApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_schemaInjection(t *testing.T) {
	type args struct {
		ctypes map[string]*meta.ChannelType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := schemaInjection(tt.args.ctypes); (err != nil) != tt.wantErr {
				t.Errorf("schemaInjection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_recursiveSchemaInjection(t *testing.T) {
	type args struct {
		apps map[string]*meta.App
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := recursiveSchemaInjection(tt.args.apps); (err != nil) != tt.wantErr {
				t.Errorf("recursiveSchemaInjection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
