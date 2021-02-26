package cli

import (
	"errors"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/controller"
	"gitlab.inspr.dev/inspr/core/pkg/controller/mocks"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gopkg.in/yaml.v2"
)

func TestNewApplyApp(t *testing.T) {
	appWithoutNameBytes, _ := yaml.Marshal(meta.App{})
	appDefaultBytes, _ := yaml.Marshal(meta.App{Meta: meta.Metadata{Name: "mock"}})
	type args struct {
		a controller.AppInterface
		b []byte
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "default_test",
			args: args{
				a: mocks.NewAppMock(nil),
				b: appDefaultBytes,
			},
			want: nil,
		},
		{
			name: "app_without_name",
			args: args{
				a: mocks.NewAppMock(nil),
				b: appWithoutNameBytes,
			},
			want: ierrors.NewError().Message("dapp without name").Build(),
		},
		{
			name: "error_testing",
			args: args{
				a: mocks.NewAppMock(errors.New("new error")),
				b: appDefaultBytes,
			},
			want: errors.New("new error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewApplyApp(tt.args.a)

			r := got(tt.args.b, nil)

			if r != nil && tt.want != nil {
				if r.Error() != tt.want.Error() {
					t.Errorf("NewApplyApp() = %v, want %v", r.Error(), tt.want.Error())
				}
			} else {
				if r != tt.want {
					t.Errorf("NewApplyApp() = %v, want %v", r, tt.want)
				}
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
		{
			name: "Valid schema injection",
			args: args{
				ctypes: map[string]*meta.ChannelType{
					"ct1": {
						Meta: meta.Metadata{
							Name: "ct1",
						},
						Schema: "test/schema_example.schema",
					},
				},
			},
			wantErr: false,
		},
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
		{
			name: "Valid schema injection",
			args: args{
				apps: getAppMap(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := recursiveSchemaInjection(tt.args.apps); (err != nil) != tt.wantErr {
				t.Errorf("recursiveSchemaInjection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getAppMap() map[string]*meta.App {
	return map[string]*meta.App{
		"app1": {
			Spec: meta.AppSpec{
				ChannelTypes: map[string]*meta.ChannelType{
					"ct1": {
						Meta: meta.Metadata{
							Name: "ct1",
						},
						Schema: "test/schema_example.schema",
					},
				},
			},
		},
		"app2": {
			Spec: meta.AppSpec{
				ChannelTypes: map[string]*meta.ChannelType{
					"ct2": {
						Meta: meta.Metadata{
							Name: "ct2",
						},
						Schema: "test/schema_example.schema",
					},
				},
				Apps: map[string]*meta.App{
					"app3": {
						Spec: meta.AppSpec{
							ChannelTypes: map[string]*meta.ChannelType{
								"ct3": {
									Meta: meta.Metadata{
										Name: "ct3",
									},
									Schema: "test/schema_example.schema",
								},
							},
						},
					},
				},
			},
		},
	}
}
