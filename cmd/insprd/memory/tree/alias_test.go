package tree

import (
	"crypto/sha256"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func TestMemoryManager_Alias(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	tests := []struct {
		name   string
		fields fields
		want   memory.AliasMemory
	}{
		{
			name: "It should return a pointer to AliasMemoryManager.",
			fields: fields{
				root: getMockAlias(),
			},
			want: &AliasMemoryManager{
				&MemoryManager{
					root: getMockAlias(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmm := &MemoryManager{
				root: tt.fields.root,
			}
			if got := tmm.Alias(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryManager.Alias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAliasMemoryManager_CreateAlias(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
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
			setTree(&MockManager{
				MemoryManager: &MemoryManager{
					root: tt.fields.root,
					tree: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			amm := GetTreeMemory().Alias()
			if err := amm.CreateAlias(tt.args.query, tt.args.targetBoundary, tt.args.alias); (err != nil) != tt.wantErr {
				t.Errorf("AliasMemoryManager.CreateAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAliasMemoryManager_GetAlias(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
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
			setTree(&MockManager{
				MemoryManager: &MemoryManager{
					root: tt.fields.root,
					tree: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			amm := GetTreeMemory().Alias()

			got, err := amm.GetAlias(tt.args.context, tt.args.aliasKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("AliasMemoryManager.GetAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AliasMemoryManager.GetAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAliasMemoryManager_UpdateAlias(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
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
			setTree(&MockManager{
				MemoryManager: &MemoryManager{
					root: tt.fields.root,
					tree: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			amm := GetTreeMemory().Alias()

			if err := amm.UpdateAlias(tt.args.context, tt.args.aliasKey, tt.args.alias); (err != nil) != tt.wantErr {
				t.Errorf("AliasMemoryManager.UpdateAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAliasMemoryManager_DeleteAlias(t *testing.T) {
	type fields struct {
		MemoryManager *MemoryManager
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amm := &AliasMemoryManager{
				MemoryManager: tt.fields.MemoryManager,
			}
			if err := amm.DeleteAlias(tt.args.context, tt.args.aliasKey); (err != nil) != tt.wantErr {
				t.Errorf("AliasMemoryManager.DeleteAlias() error = %v, wantErr %v", err, tt.wantErr)
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
		// TODO: Add test cases.
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
			SHA256:      "",
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
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"appUpdate1": {},
							"appUpdate2": {},
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
						ChannelTypes: map[string]*meta.ChannelType{
							"ctUpdate1": {
								Meta: meta.Metadata{
									Name:        "ctUpdate1",
									Reference:   "app1.ctUpdate1",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
								},
								ConnectedChannels: []string{"ch2app1Update", "ch1app1"},
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
						Type: "channelType1",
					},
				},
				"channel2": {
					Meta: meta.Metadata{
						Name:   "channel2",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "channelType1",
					},
				},
			},
			ChannelTypes: map[string]*meta.ChannelType{
				"channelType1": {
					Meta: meta.Metadata{
						Name: "channelType1",
					},
					Schema: string(sha256.New().Sum([]byte("hello"))),
				},
			},
			Boundary: meta.AppBoundary{
				Input:  []string{},
				Output: []string{},
			},
			Aliases: map[string]*meta.Alias{
				"app1.aliaschannel": {
					Target: "channel2",
				},
			},
		},
	}
	return &root
}
