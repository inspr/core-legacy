package cli

import (
	"errors"
	"io"
	"os"
	"reflect"
	"testing"

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
				t.Errorf("getFactory() = %v, want %v", got, tt.want)
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
		name          string
		fields        fields
		args          args
		checkFunction func(t *testing.T, runMethod RunMethod)
		wantErr       bool
	}{
		{
			name: "It should get the method for the given component",
			fields: fields{
				applyDict: map[meta.Component]RunMethod{
					{
						Kind:       "app",
						APIVersion: "v1",
					}: func(b []byte, out io.Writer) error {
						return errors.New("just a example to test the function return")
					},
				},
			},
			args: args{
				meta.Component{
					Kind:       "app",
					APIVersion: "v1",
				},
			},
			checkFunction: func(t *testing.T, runMethod RunMethod) {
				foo := []byte("foo")
				got := runMethod(foo, os.Stdout).Error()
				if got != "just a example to test the function return" {
					t.Errorf("ApplyFactory.GetRunMethod() = %v, want %v", got, "Just a example to test the function return")
				}
			},
			wantErr: false,
		},
		{
			name: "It should return a error - component not registered",
			fields: fields{
				applyDict: map[meta.Component]RunMethod{
					{
						Kind:       "app",
						APIVersion: "v1",
					}: func(b []byte, out io.Writer) error {
						return errors.New("just a example to test the function return")
					},
				},
			},
			args: args{
				meta.Component{
					Kind:       "invalid",
					APIVersion: "v1",
				},
			},
			wantErr: true,
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

			if tt.checkFunction != nil {
				tt.checkFunction(t, got)
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
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "It should not return a error - subscribe correctly",
			fields: fields{
				applyDict: map[meta.Component]RunMethod{},
			},
			args: args{
				component: meta.Component{
					Kind:       "app",
					APIVersion: "v1",
				},
			},
			wantErr: false,
		},
		{
			name: "It should return a error - component has invalid field",
			fields: fields{
				applyDict: map[meta.Component]RunMethod{},
			},
			args: args{
				component: meta.Component{
					Kind:       "",
					APIVersion: "v1",
				},
			},
			wantErr: true,
		},
		{
			name: "It should return a error - component already subscribed",
			fields: fields{
				applyDict: map[meta.Component]RunMethod{
					{
						Kind:       "app",
						APIVersion: "v1",
					}: func(b []byte, out io.Writer) error {
						return errors.New("just a example to test the function return")
					},
				},
			},
			args: args{
				component: meta.Component{
					Kind:       "app",
					APIVersion: "v1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			af := &ApplyFactory{
				applyDict: tt.fields.applyDict,
			}
			if err := af.Subscribe(tt.args.component, tt.args.method); (err != nil) != tt.wantErr {
				t.Errorf("applyFactory.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
