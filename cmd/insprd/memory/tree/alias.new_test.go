package tree

import (
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/meta"
)

func TestMemoryManager_AliasNew(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	tests := []struct {
		name   string
		fields fields
		want   AliasMemory
	}{
		{
			name: "It should return a pointer to AliasMemoryManagerNew.",
			fields: fields{
				root: getMockAlias(),
			},
			want: &AliasMemoryManagerNew{
				&treeMemoryManager{
					root: getMockAlias(),
				},
				logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmm := &treeMemoryManager{
				root: tt.fields.root,
			}

			if got := tmm.AliasNew(); !reflect.DeepEqual(got.(*AliasMemoryManagerNew).root, tt.want.(*AliasMemoryManagerNew).root) {
				t.Errorf("MemoryManager.Alias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAliasMemoryManagerNew_Create(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	type args struct {
		scope string
		alias *meta.Alias
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "app doesn't exist, it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "invalid.app",
				alias: &meta.Alias{},
			},
			wantErr: true,
		},
		{
			name: "alias already exist - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				alias: &meta.Alias{
					Meta: meta.Metadata{
						Name: "my_alias",
					},
					Source:      "appUpdate1",
					Resource:    "route_1",
					Destination: "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid source - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				alias: &meta.Alias{
					Meta: meta.Metadata{
						Name: "my_new_alias",
					},
					Source: "invalid_child",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid resource - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				alias: &meta.Alias{
					Meta: meta.Metadata{
						Name: "my_new_alias",
					},
					Source:      "appUpdate1",
					Resource:    "invalid",
					Destination: "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid destination - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				alias: &meta.Alias{
					Meta: meta.Metadata{
						Name: "my_new_alias",
					},
					Source:      "appUpdate1",
					Resource:    "route_1",
					Destination: "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "valid creation - it should not return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				alias: &meta.Alias{
					Meta: meta.Metadata{
						Name: "my_new_alias",
					},
					Source:      "appUpdate1",
					Resource:    "route_1",
					Destination: "appUpdate2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &treeMemoryManager{
				root: tt.fields.root,
				tree: tt.fields.root,
			}
			amm := mem.Alias()
			if err := amm.Create(tt.args.scope, tt.args.alias); (err != nil) != tt.wantErr {
				t.Errorf("AliasMemoryManager.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAliasMemoryManagerNew_Get(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	type args struct {
		scope string
		name  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *meta.Alias
		wantErr bool
	}{
		{
			name: "invalid scope - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "invalid.app",
				name:  "my_alias",
			},
			wantErr: true,
		},
		{
			name: "alias name not exist - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				name:  "invalid.alias",
			},
			wantErr: true,
		},
		{
			name: "Valid scope, alias name exist - it should not return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				name:  "my_alias",
			},
			wantErr: false,
			want: &meta.Alias{
				Meta: meta.Metadata{
					Name: "my_alias",
				},
				Resource:    "channel1",
				Source:      "",
				Destination: "appUpdate1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &treeMemoryManager{
				root: tt.fields.root,
				tree: tt.fields.root,
			}
			amm := mem.Alias()

			got, err := amm.Get(tt.args.scope, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("AliasMemoryManager.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AliasMemoryManager.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAliasMemoryManagerNew_Update(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	type args struct {
		scope string
		alias *meta.Alias
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "invalid scope - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "invalid.app",
				alias: &meta.Alias{},
			},
			wantErr: true,
		},
		{
			name: "alias do not exist - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				alias: &meta.Alias{
					Meta: meta.Metadata{
						Name: "invalid.alias",
					},
					Source: "invalid_child",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid source - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				alias: &meta.Alias{
					Meta: meta.Metadata{
						Name: "my_alias",
					},
					Source: "invalid_child",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid resource - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				alias: &meta.Alias{
					Meta: meta.Metadata{
						Name: "my_alias",
					},
					Source:      "appUpdate1",
					Resource:    "invalid",
					Destination: "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid destination - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				alias: &meta.Alias{
					Meta: meta.Metadata{
						Name: "my_alias",
					},
					Source:      "appUpdate1",
					Resource:    "route_1",
					Destination: "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "Valid update - it should not return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				alias: &meta.Alias{
					Meta: meta.Metadata{
						Name: "my_alias",
					},
					Source:      "appUpdate1",
					Resource:    "route_1",
					Destination: "appUpdate2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &treeMemoryManager{
				root: tt.fields.root,
				tree: tt.fields.root,
			}
			amm := mem.Alias()

			if err := amm.Update(tt.args.scope, tt.args.alias); (err != nil) != tt.wantErr {
				t.Errorf("AliasMemoryManager.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
