package cli

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func TestGetFactory(t *testing.T) {
	tests := []struct {
		name string
		want *ApplyFactory
	}{
		{
			name: "It should create and return a ApplyFactory",
			want: &ApplyFactory{
				applyDict: map[meta.Component]RunMethod{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFactory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyFactory_GetRunMethod(t *testing.T) {
	type fields struct {
		applyDict map[meta.Component]RunMethod
	}
	type args struct {
		component meta.Component
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    RunMethod
		wantErr bool
	}{
		{
			name: "It should get the method for the given component",
			fields: fields{
				applyDict: map[meta.Component]RunMethod{
					{
						Kind:       "app",
						APIVersion: "v1",
					}: func(c *cobra.Command, s []string) {
						fmt.Println("It runs the function")
					},
				},
			},
			args: args{
				meta.Component{
					Kind:       "app",
					APIVersion: "v1",
				},
			},
			want: func(c *cobra.Command, s []string) {
				fmt.Println("It runs the function")
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			af := &ApplyFactory{
				applyDict: tt.fields.applyDict,
			}
			got, err := af.GetRunMethod(tt.args.component)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyFactory.GetRunMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyFactory.GetRunMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyFactory_Subscribe(t *testing.T) {
	type fields struct {
		applyDict map[meta.Component]RunMethod
	}
	type args struct {
		component meta.Component
		method    RunMethod
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			af := &ApplyFactory{
				applyDict: tt.fields.applyDict,
			}
			af.Subscribe(tt.args.component, tt.args.method)
		})
	}
}
