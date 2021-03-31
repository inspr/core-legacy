package tree

import (
	"crypto/sha256"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

func TestMemoryManager_Channels(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	tests := []struct {
		name   string
		fields fields
		want   memory.ChannelMemory
	}{
		{
			name: "It should return a pointer to ChannelMemoryManager.",
			fields: fields{
				root: getMockChannels(),
			},
			want: &ChannelMemoryManager{
				&MemoryManager{
					root: getMockChannels(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmm := &MemoryManager{
				root: tt.fields.root,
			}
			if got := tmm.Channels(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryManager.Channels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelMemoryManager_GetChannel(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		context string
		chName  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *meta.Channel
		wantErr bool
	}{
		{
			name: "It should return a valid Channel",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "",
				chName:  "channel1",
			},
			wantErr: false,
			want: &meta.Channel{
				Meta: meta.Metadata{
					Name:   "channel1",
					Parent: "",
				},
				ConnectedApps: []string{"app1"},
				Spec: meta.ChannelSpec{
					Type: "channelType1",
				},
			},
		},
		{
			name: "It should return a invalid Channel on a valid App",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "",
				chName:  "channel3",
			},
			wantErr: true,
		},
		{
			name: "It should return a invalid Channel on a invalid App",
			fields: fields{
				root:   getMockChannels(),
				appErr: ierrors.NewError().NotFound().Build(),
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "invalid.context",
				chName:  "channel1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTree(&MockManager{
				MemoryManager: &MemoryManager{
					root: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			chh := GetTreeMemory().Channels()
			got, err := chh.Get(tt.args.context, tt.args.chName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChannelMemoryManager.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelMemoryManager.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelMemoryManager_CreateChannel(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		context string
		ch      *meta.Channel
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		want          *meta.Channel
		checkFunction func() (bool, string)
	}{
		{
			name: "It should create a new Channel on a valid App",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "",
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name:   "channel3",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "channelType1",
					},
				},
			},
			wantErr: false,
			want: &meta.Channel{
				Meta: meta.Metadata{
					Name:   "channel3",
					Parent: "",
				},
				Spec: meta.ChannelSpec{
					Type: "channelType1",
				},
			},
		},
		{
			name: "It should not create a new Channel because it already exists",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "",
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name:   "channel1",
						Parent: "",
					},
					Spec: meta.ChannelSpec{},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "It should not create a new Channel because the context is invalid",
			fields: fields{
				root:   getMockChannels(),
				appErr: ierrors.NewError().NotFound().Build(),
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "invalid.context",
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name:   "channel3",
						Parent: "",
					},
					Spec: meta.ChannelSpec{},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "It should not create a channel because the channelType is invalid",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "",
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name:   "channel3",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "ct1",
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "It should create a channel because the channelType is valid",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: false,
			},
			args: args{
				context: "",
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name:   "channel3",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "channelType1",
					},
				},
			},
			wantErr: false,
			want:    nil,
			checkFunction: func() (bool, string) {
				am := GetTreeMemory().ChannelTypes()
				ct, err := am.Get("", "channelType1")
				if err != nil {
					return false, "cant get channelType 'channelType1'"
				}

				if !utils.Includes(ct.ConnectedChannels, "channel3") {
					return false, "connectedChannels of channelType1 dont have channel3"
				}
				return true, ""
			},
		},
		{
			name: "Invalid channel name - doesn't create channel",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "",
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name:   "channel3/",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "channelType1",
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTree(&MockManager{
				MemoryManager: &MemoryManager{
					root: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			chh := GetTreeMemory().Channels()
			err := chh.CreateChannel(tt.args.context, tt.args.ch)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChannelMemoryManager.CreateChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {

				got, err := chh.Get(tt.args.context, tt.want.Meta.Name)
				if !tt.wantErr {
					if !ValidateUUID(got.Meta.UUID) {
						t.Errorf("ChannelMemoryManager.Create() invalid UUID, uuid=%v", got.Meta.UUID)
					}
				}
				if (err != nil) || !CompareWithoutUUID(got, tt.want) {
					t.Errorf("ChannelMemoryManager.Get() = %v, want %v", got, tt.want)
				}
			}
			if tt.checkFunction != nil {
				if passed, msg := tt.checkFunction(); !passed {
					t.Errorf("check function not passed: " + msg)
				}
			}
		})
	}
}

func TestChannelMemoryManager_DeleteChannel(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		context string
		chName  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *meta.Channel
	}{
		{
			name: "It should delete a Channel on a valid App",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "",
				chName:  "channel2",
			},
			wantErr: false,
			want:    nil,
		},
		{
			name: "It should not delete the channel, because it does not exist",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				chName:  "channel3",
				context: "",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "It shoud not delete the Channel because the context is invalid.",
			fields: fields{
				root:   nil,
				appErr: ierrors.NewError().NotFound().Build(),
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "invalid.context",
				chName:  "channel1",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "It should not delete the Channel because it's been used by an app.",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "",
				chName:  "channel1",
			},
			wantErr: true,
			want: &meta.Channel{
				Meta: meta.Metadata{
					Name:   "channel1",
					Parent: "",
				},
				ConnectedApps: []string{"app1"},
				Spec: meta.ChannelSpec{
					Type: "channelType1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTree(&MockManager{
				MemoryManager: &MemoryManager{
					root: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			chh := GetTreeMemory().Channels()
			if err := chh.DeleteChannel(tt.args.context, tt.args.chName); (err != nil) != tt.wantErr {
				t.Errorf("ChannelMemoryManager.DeleteChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, _ := chh.Get(tt.args.context, tt.args.chName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelMemoryManager.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelMemoryManager_UpdateChannel(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		context string
		ch      *meta.Channel
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *meta.Channel
	}{
		{
			name: "It should update a Channel on a valid App",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "",
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name: "channel1",
						Annotations: map[string]string{
							"update1": "update1",
						},
						Parent: "",
					},
					ConnectedApps: []string{"app1"},
					Spec: meta.ChannelSpec{
						Type: "channelType1",
					},
				},
			},
			wantErr: false,
			want: &meta.Channel{
				Meta: meta.Metadata{
					Name: "channel1",
					Annotations: map[string]string{
						"update1": "update1",
					},
					Parent: "",
				},
				ConnectedApps: []string{"app1"},
				Spec: meta.ChannelSpec{
					Type: "channelType1",
				},
			},
		},
		{
			name: "It should not update a Channel because it does not exist",
			fields: fields{
				root:   getMockChannels(),
				appErr: nil,
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "",
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name: "channel3",
						Annotations: map[string]string{
							"update1": "update1",
						},
						Parent: "",
					},
					Spec: meta.ChannelSpec{},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "It should not update a Channel because the context is invalid",
			fields: fields{
				root:   getMockChannels(),
				appErr: ierrors.NewError().NotFound().Build(),
				mockA:  true,
				mockC:  false,
				mockCT: true,
			},
			args: args{
				context: "invalid.context",
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name: "channel1",
						Annotations: map[string]string{
							"update1": "update1",
						},
						Parent: "",
					},
					Spec: meta.ChannelSpec{},
				},
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTree(&MockManager{
				MemoryManager: &MemoryManager{
					root: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			chh := GetTreeMemory().Channels()
			if err := chh.UpdateChannel(tt.args.context, tt.args.ch); (err != nil) != tt.wantErr {
				t.Errorf("ChannelMemoryManager.UpdateChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				got, err := chh.Get(tt.args.context, tt.want.Meta.Name)

				if (err != nil) || !CompareWithUUID(got, tt.want) {
					t.Errorf("ChannelMemoryManager.Get() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func getMockChannels() *meta.App {
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
									UUID:        "",
								},
								ConnectedChannels: []string{"ch2app1Update", "ch1app1"},
							},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{"channel1"},
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
		},
	}
	return &root
}
