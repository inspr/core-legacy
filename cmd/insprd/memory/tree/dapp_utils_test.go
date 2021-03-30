package tree

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

func Test_validAppStructure(t *testing.T) {
	type args struct {
		app       meta.App
		parentApp meta.App
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "All valid structures",
			args: args{
				app: meta.App{
					Meta: meta.Metadata{
						Name:        "app5",
						Reference:   "app2.app5",
						Annotations: map[string]string{},
						Parent:      "app2",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "nodeApp5",
								Reference:   "app5.nodeApp5",
								Annotations: map[string]string{},
								Parent:      "app2",
								SHA256:      "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeApp5",
							},
						},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2"},
							Output: []string{"ch2app2"},
						},
					},
				},
				parentApp: *getMockApp().Spec.Apps["app2"],
			},
			want: "",
		},
		{
			name: "invalidapp name - empty",
			args: args{
				app: meta.App{
					Meta: meta.Metadata{
						Name:        "",
						Reference:   "app2.app4",
						Annotations: map[string]string{},
						Parent:      "app2",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "nodeApp4",
								Reference:   "app4.nodeApp4",
								Annotations: map[string]string{},
								Parent:      "app3",
								SHA256:      "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeApp3",
							},
						},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2"},
							Output: []string{"ch2app2"},
						},
					},
				},
				parentApp: *getMockApp().Spec.Apps["app2"],
			},
			want: "invalid dApp name;",
		},
		{
			name: "invalidapp substructure",
			args: args{
				app: meta.App{
					Meta: meta.Metadata{
						Name:        "app5",
						Reference:   "app2.app5",
						Annotations: map[string]string{},
						Parent:      "app2",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "nodeApp5",
								Reference:   "app5.nodeApp5",
								Annotations: map[string]string{},
								Parent:      "app5",
								SHA256:      "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeApp5",
							},
						},
						Apps: map[string]*meta.App{
							"invalidApp": {},
						},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2"},
							Output: []string{"ch2app2"},
						},
					},
				},
				parentApp: *getMockApp().Spec.Apps["app2"],
			},
			want: "invalid substructure;",
		},
		{
			name: "invalidapp - parent has Node structure",
			args: args{
				app: meta.App{
					Meta: meta.Metadata{
						Name:        "app4",
						Reference:   "app2.app4",
						Annotations: map[string]string{},
						Parent:      "app2",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "nodeApp4",
								Reference:   "app4.nodeApp4",
								Annotations: map[string]string{},
								Parent:      "app3",
								SHA256:      "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeApp3",
							},
						},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app1"},
							Output: []string{"ch2app1"},
						},
					},
				},
				parentApp: *getMockApp().Spec.Apps["appNode"],
			},
			want: "parent has Node;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validAppStructure(&tt.args.app, &tt.args.parentApp); got != tt.want {
				t.Errorf("validAppStructure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeIsEmpty(t *testing.T) {
	type args struct {
		node meta.Node
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Node is empty",
			args: args{
				node: meta.Node{},
			},
			want: true,
		},
		{
			name: "Node isn't empty",
			args: args{
				node: meta.Node{
					Meta: meta.Metadata{
						Name:        "nodeApp1",
						Reference:   "app1.nodeApp1",
						Annotations: map[string]string{},
						Parent:      "app1",
						SHA256:      "",
					},
					Spec: meta.NodeSpec{
						Image: "imageNodeApp1",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nodeIsEmpty(tt.args.node); got != tt.want {
				t.Errorf("nodeIsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validBoundaries(t *testing.T) {
	type args struct {
		appName        string
		bound          meta.AppBoundary
		parentChannels map[string]*meta.Channel
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid boundary",
			args: args{
				appName: "thenewapp",
				bound: meta.AppBoundary{
					Input:  []string{"ch1app1"},
					Output: []string{},
				},
				parentChannels: getMockApp().Spec.Apps["app1"].Spec.Channels,
			},
			want: "",
		},
		{
			name: "invalidboundary - parent without channels",
			args: args{
				appName: "",
				bound: meta.AppBoundary{
					Input:  []string{"ch1app2"},
					Output: []string{},
				},
				parentChannels: getMockApp().Spec.Apps["app2"].Spec.Apps["app3"].Spec.Channels,
			},
			want: "invalid app boundary - channel 'ch1app2' doesnt exist in parent app;",
		},
		{
			name: "invalid input boundary",
			args: args{
				appName: "app3",
				bound: meta.AppBoundary{
					Input:  []string{"ch1app1"},
					Output: []string{},
				},
				parentChannels: getMockApp().Spec.Apps["app2"].Spec.Channels,
			},
			want: "invalid app boundary - channel 'ch1app1' doesnt exist in parent app;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validBoundaries(tt.args.appName, tt.args.bound, tt.args.parentChannels); got != tt.want {
				t.Errorf("validBoundaries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getParentApp(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockC  bool
		mockCT bool
		mockA  bool
	}
	type args struct {
		sonQuery string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *meta.App
		wantErr bool
	}{
		{
			name: "Parent is the root",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				sonQuery: "app1",
			},
			wantErr: false,
			want:    getMockApp(),
		},
		{
			name: "Parent is another app",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				sonQuery: "app2.app3",
			},
			wantErr: false,
			want:    getMockApp().Spec.Apps["app2"],
		},
		{
			name: "invalidquery",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				sonQuery: "invalid.query",
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
					tree: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			got, err := getParentApp(tt.args.sonQuery)
			if (err != nil) != tt.wantErr {
				t.Errorf("getParentApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getParentApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkAndUpdates(t *testing.T) {
	type args struct {
		app *meta.App
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{
			name: "valid channel structure - it shouldn't return a error",
			args: args{
				app: &meta.App{
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
							"thenewapp": {
								Meta: meta.Metadata{
									Name:        "thenewapp",
									Reference:   "app1.thenewapp",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1app1"},
										Output: []string{},
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
								ConnectedApps: []string{"thenewapp"},
								Spec: meta.ChannelSpec{
									Type: "newChannelType",
								},
							},
							"ch2app1": {
								Meta: meta.Metadata{
									Name:   "ch2app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{
							"newChannelType": {
								Meta: meta.Metadata{
									Name:        "newChannelType",
									Reference:   "app1.newChannelType",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
								},
							},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			want:  true,
			want1: "",
		},
		{
			name: "invalid channel: using non-existent channel type",
			args: args{
				app: &meta.App{
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
							"thenewapp": {
								Meta: meta.Metadata{
									Name:        "thenewapp",
									Reference:   "app1.thenewapp",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1app1"},
										Output: []string{},
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
								ConnectedApps: []string{"thenewapp"},
								Spec: meta.ChannelSpec{
									Type: "invalidType",
								},
							},
							"ch2app1": {
								Meta: meta.Metadata{
									Name:   "ch2app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{
							"newChannelType": {
								Meta: meta.Metadata{
									Name:        "newChannelType",
									Reference:   "app1.newChannelType",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
								},
							},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			want:  false,
			want1: "invalid channel: using non-existent channel type;",
		},
		{
			name: "invalid channel structure - it should return a name channel error",
			args: args{
				app: &meta.App{
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
							"thenewapp": {
								Meta: meta.Metadata{
									Name:        "thenewapp",
									Reference:   "app1.thenewapp",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1app1"},
										Output: []string{},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"invalid.channel.name": {
								Meta: meta.Metadata{
									Name:   "ch1app1",
									Parent: "",
								},
								ConnectedApps: []string{"thenewapp"},
								Spec: meta.ChannelSpec{
									Type: "newChannelType",
								},
							},
							"ch2app1": {
								Meta: meta.Metadata{
									Name:   "ch2app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{
							"newChannelType": {
								Meta: meta.Metadata{
									Name:        "newChannelType",
									Reference:   "app1.newChannelType",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
								},
							},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			want:  false,
			want1: "invalid channel name: invalid.channel.name",
		},
		{
			name: "valid channel structure - it shouldn't return a error",
			args: args{
				app: &meta.App{
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
							"thenewapp": {
								Meta: meta.Metadata{
									Name:        "thenewapp",
									Reference:   "app1.thenewapp",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1app1"},
										Output: []string{},
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
								ConnectedApps: []string{"thenewapp"},
								Spec: meta.ChannelSpec{
									Type: "newChannelType",
								},
							},
							"ch2app1": {
								Meta: meta.Metadata{
									Name:   "ch2app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{
							"invalid.channel.type": {
								Meta: meta.Metadata{
									Name:        "newChannelType",
									Reference:   "app1.newChannelType",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
								},
							},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			want:  false,
			want1: "invalid channelType name: invalid.channel.type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := checkAndUpdates(tt.args.app)
			if got != tt.want {
				t.Errorf("checkChannels() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("checkChannels() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAppMemoryManager_connectAppBoundary(t *testing.T) {
	type fields struct {
		MemoryManager *MemoryManager
		root          *meta.App
		appErr        error
		mockA         bool
		mockC         bool
		mockCT        bool
	}
	type args struct {
		app *meta.App
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantedApps  map[string]utils.StringArray
		wantedAlias map[string]utils.StringArray
		wantErr     bool
		sourceApp   string
	}{
		{
			name: "Valid - direct connect",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				app: getMockApp().Spec.Apps["bound"].Spec.Apps["bound6"],
			},
			wantedApps: map[string]utils.StringArray{
				"bdch1": {"bound6"},
				"bdch2": nil,
			},
			wantedAlias: map[string]utils.StringArray{
				"bdch1": {"bound2.alias1"},
				"bdch2": {"bound2.alias2", "bound4.alias3", "bound6.alias3"},
			},
			sourceApp: "bound",
			wantErr:   false,
		},
		{
			name: "Invalid - app with bad parent",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				app: getMockApp().Spec.Apps["bound"].Spec.Apps["boundNP"],
			},
			wantErr: true,
		},
		{
			name: "Invalid - parent with bad alias",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				app: getMockApp().Spec.Apps["bound"].Spec.Apps["bound6"].Spec.Apps["bound7"],
			},
			wantErr: true,
		},
		{
			name: "Valid - resolve boundary through recursive alias",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				app: getMockApp().Spec.Apps["bound"].Spec.Apps["bound2"].Spec.Apps["bound3"],
			},
			sourceApp: "bound2",
			wantErr:   false,
		},
		{
			name: "Invalid - bad reference",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				app: getMockApp().Spec.Apps["bound"].Spec.Apps["bound5"],
			},
			wantErr: true,
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
			amm := GetTreeMemory().Apps().(*AppMemoryManager)
			err := amm.connectAppBoundary(tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.connectAppsThroughAliases() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, _ := amm.Get(tt.sourceApp)
			for ch, conn := range tt.wantedApps {
				if len(got.Spec.Channels[ch].ConnectedApps) != len(conn) {
					t.Errorf("AppMemoryManager.connectAppBoundary() on %s.ConnectedApps = %v, want = %v", ch, got.Spec.Channels[ch].ConnectedApps, conn)
					return
				}
				for _, app := range conn {
					if !got.Spec.Channels[ch].ConnectedApps.Contains(app) {
						t.Errorf("AppMemoryManager.connectAppBoundary() on %s.ConnectedApps = %v, want = %v", ch, got.Spec.Channels[ch].ConnectedApps, conn)
						return
					}
				}
			}
			for ch, conn := range tt.wantedAlias {
				if len(got.Spec.Channels[ch].ConnectedAliases) != len(conn) {
					t.Errorf("AppMemoryManager.connectAppBoundary()  on %s.ConnectedAliases = %v, want = %v", ch, got.Spec.Channels[ch].ConnectedAliases, conn)
					return
				}
				for _, alias := range conn {
					if !got.Spec.Channels[ch].ConnectedAliases.Contains(alias) {
						t.Errorf("AppMemoryManager.connectAppBoundary() on %s.ConnectedApps = %v, want = %v", ch, got.Spec.Channels[ch].ConnectedApps, conn)
						return
					}
				}
			}
		})
	}
}

func TestAppMemoryManager_connectAppsBoundaries(t *testing.T) {
	type fields struct {
		MemoryManager *MemoryManager
		root          *meta.App
		appErr        error
		mockA         bool
		mockC         bool
		mockCT        bool
	}
	type args struct {
		app *meta.App
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Full coverage",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				app: getMockApp().Spec.Apps["bound"].Spec.Apps["bound2"],
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
			amm := GetTreeMemory().Apps().(*AppMemoryManager)
			if err := amm.connectAppsBoundaries(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.connectAppsBoundaries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
