package tree

import (
	"crypto/sha256"
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/utils"
)

func TestMemoryManager_Alias(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	tests := []struct {
		name   string
		fields fields
		want   AliasMemory
	}{
		{
			name: "It should return a pointer to AliasMemoryManager.",
			fields: fields{
				root: getMockAlias(),
			},
			want: &AliasMemoryManager{
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

			if got := tmm.Alias(); !reflect.DeepEqual(got.(*AliasMemoryManager).root, tt.want.(*AliasMemoryManager).root) {
				t.Errorf("MemoryManager.Alias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAliasMemoryManager_Create(t *testing.T) {
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

func TestAliasMemoryManager_Get(t *testing.T) {
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

func TestAliasMemoryManager_Update(t *testing.T) {
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

func TestAliasMemoryManager_Delete(t *testing.T) {
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
			name: "alias do not exist - it should return an error",
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
			name: "alias exist but its being used by the child dapp in another alias - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				name:  "my_alias",
			},
			wantErr: true,
		},
		{
			name: "alias exist but its being used by the child dapp in its boudaries - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				name:  "my_other_alias",
			},
			wantErr: true,
		},
		{
			name: "alias exist but its being used by the parent dapp - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1",
				name:  "my_awesome_alias",
			},
			wantErr: true,
		},
		{
			name: "Valid delete - it should not return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				scope: "app1.appUpdate1",
				name:  "my_brand_new_alias",
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
			if err := amm.Delete(tt.args.scope, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("AliasMemoryManager.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getMockAlias() *meta.App {
	root := meta.App{
		Meta: meta.Metadata{
			Name:        "",
			Reference:   "",
			Annotations: map[string]string{},
			Parent:      "",
			UUID:        "",
		},
		Spec: meta.AppSpec{
			Node: meta.Node{},
			Apps: map[string]*meta.App{
				"app1": {
					Meta: meta.Metadata{
						Name:        "app1",
						Reference:   "app1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"appUpdate1": {
								Meta: meta.Metadata{
									Name: "appUpdate1",
								},
								Spec: meta.AppSpec{
									Routes: map[string]*meta.RouteConnection{
										"route_1": {},
									},
									Apps: map[string]*meta.App{
										"app_1_1": {},
									},
									Aliases: map[string]*meta.Alias{
										"my_brand_new_alias": {
											Meta: meta.Metadata{
												Name: "my_brand_new_alias",
											},
											Resource:    "my_alias",
											Destination: "app_1_1",
											Source:      "",
										},
									},
								},
							},
							"appUpdate2": {
								Spec: meta.AppSpec{
									Boundary: meta.AppBoundary{
										Output: utils.StringArray{
											"my_other_alias",
										},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"ch1app1": {
								Meta: meta.Metadata{
									Name:   "ch1app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{
									Type: "ctUpdate1",
								},
							},
							"ch2app1Update": {
								Meta: meta.Metadata{
									Name:   "ch2app1Update",
									Parent: "app1",
								},
								Spec: meta.ChannelSpec{
									Type: "ctUpdate1",
								},
							},
						},
						Types: map[string]*meta.Type{
							"ctUpdate1": {
								Meta: meta.Metadata{
									Name:        "ctUpdate1",
									Reference:   "app1.ctUpdate1",
									Annotations: map[string]string{},
									Parent:      "app1",
									UUID:        "",
								},
								ConnectedChannels: []string{"ch2app1Update", "ch1app1"},
							},
						},

						Aliases: map[string]*meta.Alias{
							"my_alias": {
								Meta: meta.Metadata{
									Name: "my_alias",
								},
								Resource:    "channel1",
								Source:      "",
								Destination: "appUpdate1",
							},
							"my_other_alias": {
								Meta: meta.Metadata{
									Name: "my_other_alias",
								},
								Resource:    "channel1",
								Source:      "",
								Destination: "appUpdate2",
							},
							"my_awesome_alias": {
								Meta: meta.Metadata{
									Name: "my_awesome_alias",
								},
								Resource:    "route_1",
								Source:      "appUpdate1",
								Destination: "",
							},
						},

						Boundary: meta.AppBoundary{
							Input:  []string{"channel1", "aliaschannel", "aliaschannel2"},
							Output: []string{},
						},
					},
				},
				"app2": {},
			},
			Channels: map[string]*meta.Channel{
				"channel1": {
					Meta: meta.Metadata{
						Name:   "channel1",
						Parent: "",
					},
					ConnectedApps: []string{"app1"},
					Spec: meta.ChannelSpec{
						Type: "type1",
					},
				},
				"channel2": {
					Meta: meta.Metadata{
						Name:   "channel2",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "type1",
					},
				},
			},
			Types: map[string]*meta.Type{
				"type1": {
					Meta: meta.Metadata{
						Name: "type1",
					},
					Schema: string(sha256.New().Sum([]byte("hello"))),
				},
			},
			Boundary: meta.AppBoundary{
				Input:  []string{"somechannel"},
				Output: []string{},
			},
			Aliases: map[string]*meta.Alias{
				"app1.aliaschannel": {
					Target: "channel2",
				},
				"app2.aliaschannel": {
					Target: "channel2",
				},

				"my_crazy_alias": {
					Meta: meta.Metadata{
						Name: "my_crazy_alias",
					},
					Resource: "my_awesome_alias",
					Source:   "app1",
				},
			},
		},
	}
	return &root
}
