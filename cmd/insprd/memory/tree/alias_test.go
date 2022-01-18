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

			if got := tmm.Alias(); !reflect.DeepEqual(got.(*AliasMemoryManagerNew).root, tt.want.(*AliasMemoryManagerNew).root) {
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
		query          string
		targetBoundary string
		alias          *meta.Alias
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
				query:          "invalid.app",
				targetBoundary: "ch1",
				alias:          &meta.Alias{},
			},
			wantErr: true,
		},
		{
			name: "app exist, but target boundary dont - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				query:          "app1",
				targetBoundary: "invalid",
				alias:          &meta.Alias{},
			},
			wantErr: true,
		},
		{
			name: "target channel don't exist in parent - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				query:          "app1",
				targetBoundary: "channel1",
				alias: &meta.Alias{
					Target: "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "alias already exist - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				query:          "app1",
				targetBoundary: "aliaschannel",
				alias: &meta.Alias{
					Target: "channel1",
				},
			},
			wantErr: true,
		},
		{
			name: "Valid query - it should not return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				query:          "app1",
				targetBoundary: "aliaschannel2",
				alias: &meta.Alias{
					Target: "channel1",
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
			if err := amm.Create(tt.args.query, tt.args.alias); (err != nil) != tt.wantErr {
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
		context  string
		aliasKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *meta.Alias
		wantErr bool
	}{
		{
			name: "invalid context - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				context:  "invalid.context",
				aliasKey: "app1.aliaschannel",
			},
			wantErr: true,
		},
		{
			name: "alias key not exist - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				context:  "",
				aliasKey: "invalid.alias",
			},
			wantErr: true,
		},
		{
			name: "Valid query, alias key exist - it should not return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				context:  "",
				aliasKey: "app1.aliaschannel",
			},
			wantErr: false,
			want: &meta.Alias{
				Target: "channel2",
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

			got, err := amm.Get(tt.args.context, tt.args.aliasKey)
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
		context  string
		aliasKey string
		alias    *meta.Alias
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "invalid context - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				context:  "invalid.context",
				aliasKey: "app1.aliaschannel",
				alias:    &meta.Alias{},
			},
			wantErr: true,
		},
		{
			name: "alias key not exist - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				context:  "",
				aliasKey: "invalid.alias",
				alias:    &meta.Alias{},
			},
			wantErr: true,
		},
		{
			name: "target channel don't exist in parent - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				context:  "",
				aliasKey: "app1.aliaschannel",
				alias: &meta.Alias{
					Target: "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "Valid query - it should not return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				context:  "",
				aliasKey: "app1.aliaschannel",
				alias: &meta.Alias{
					Target: "channel1",
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

			if err := amm.Update(tt.args.context, tt.args.alias); (err != nil) != tt.wantErr {
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
		context  string
		aliasKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "invalid context - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				context:  "invalid.context",
				aliasKey: "app1.aliaschannel",
			},
			wantErr: true,
		},
		{
			name: "alias key not exist - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				context:  "",
				aliasKey: "invalid.alias",
			},
			wantErr: true,
		},
		{
			name: "alias exist but its being used - it should return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				context:  "",
				aliasKey: "app1.aliaschannel",
			},
			wantErr: true,
		},
		{
			name: "Valid query - it should not return an error",
			fields: fields{
				root: getMockAlias(),
			},
			args: args{
				context:  "",
				aliasKey: "app2.aliaschannel",
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
			if err := amm.Delete(tt.args.context, tt.args.aliasKey); (err != nil) != tt.wantErr {
				t.Errorf("AliasMemoryManager.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validTargetChannel(t *testing.T) {
	type args struct {
		parentApp     *meta.App
		targetChannel string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "target channel doesn't exist in app - it should return an error",
			args: args{
				parentApp:     getMockAlias(),
				targetChannel: "invalid_channel",
			},
			wantErr: true,
		},
		{
			name: "channel exist in app channels - it should not return an error",
			args: args{
				parentApp:     getMockAlias(),
				targetChannel: "channel1",
			},
			wantErr: false,
		},
		{
			name: "channel exist in app boundary - it should not return an error",
			args: args{
				parentApp:     getMockAlias(),
				targetChannel: "somechannel",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validTargetChannel(tt.args.parentApp, tt.args.targetChannel); (err != nil) != tt.wantErr {
				t.Errorf("validTargetChannel() error = %v, wantErr %v", err, tt.wantErr)
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
										"route_1": &meta.RouteConnection{},
									},
									Apps: map[string]*meta.App{
										"app_1_1": &meta.App{},
									},
									Aliases: map[string]*meta.Alias{
										"my_brand_new_alias": &meta.Alias{
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
