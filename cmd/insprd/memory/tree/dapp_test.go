package tree

import (
	"fmt"
	"testing"

	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/cmd/sidecars"
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/meta"
	metautils "inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/meta/utils/diff"
	"inspr.dev/inspr/pkg/utils"
)

func getMockApp() *meta.App {
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
				"appNode": {
					Meta: meta.Metadata{
						Name:        "appNode",
						Reference:   "appNode",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "appNode",
								Reference:   "appNode.appNode",
								Annotations: map[string]string{},
								Parent:      "appNode",
								UUID:        "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeAppNode",
							},
						},
						Apps: map[string]*meta.App{},
						Channels: map[string]*meta.Channel{
							"ch1app1": {
								Meta: meta.Metadata{
									Name:   "ch1app1",
									Parent: "",
								},
								ConnectedApps: []string{"thenewapp"},
								Spec:          meta.ChannelSpec{},
							},
							"ch2app1": {
								Meta: meta.Metadata{
									Name:   "ch2app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						Types: map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
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
							"thenewapp": {
								Meta: meta.Metadata{
									Name:        "thenewapp",
									Reference:   "app1.thenewapp",
									Annotations: map[string]string{},
									Parent:      "app1",
									UUID:        "",
								},
								Spec: meta.AppSpec{
									Apps:     map[string]*meta.App{},
									Channels: map[string]*meta.Channel{},
									Types:    map[string]*meta.Type{},
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
								Spec:          meta.ChannelSpec{},
							},
							"ch2app1": {
								Meta: meta.Metadata{
									Name:   "ch2app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						Types: map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
				"app2": {
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "app2",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"app3": {
								Meta: meta.Metadata{
									Name:        "app3",
									Reference:   "app2.app3",
									Annotations: map[string]string{},
									Parent:      "app2",
									UUID:        "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "app3",
											Reference:   "app3.nodeApp2",
											Annotations: map[string]string{},
											Parent:      "app3",
											UUID:        "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeApp3",
										},
									},
									Apps:     map[string]*meta.App{},
									Channels: map[string]*meta.Channel{},
									Types:    map[string]*meta.Type{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1app2"},
										Output: []string{"ch2app2"},
									},
								},
							},
							"app4": {
								Meta: meta.Metadata{
									Name:        "app4",
									Reference:   "app2.app4",
									Annotations: map[string]string{},
									Parent:      "app2",
									UUID:        "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "app4",
											Reference:   "app4.nodeApp4",
											Annotations: map[string]string{},
											Parent:      "app2.app4",
											UUID:        "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeApp4",
										},
									},
									Apps: map[string]*meta.App{},
									Channels: map[string]*meta.Channel{
										"ch1app4": {
											Meta: meta.Metadata{
												Name:   "ch1app4",
												Parent: "app4",
											},
											Spec: meta.ChannelSpec{
												Type: "ctapp4",
											},
										},
										"ch2app4": {
											Meta: meta.Metadata{
												Name:   "ch2app4",
												Parent: "",
											},
											Spec: meta.ChannelSpec{},
										},
									},
									Types: map[string]*meta.Type{
										"ctapp4": {
											Meta: meta.Metadata{
												Name:        "ctUpdate1",
												Reference:   "app1.ctUpdate1",
												Annotations: map[string]string{},
												Parent:      "app1",
												UUID:        "",
											},
										},
									},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1app2"},
										Output: []string{"ch2app2"},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"ch1app2": {
								Meta: meta.Metadata{
									Name:   "ch1app2",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
							"ch2app2": {
								Meta: meta.Metadata{
									Name:   "ch2app2",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						Types: map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
				"bound": {
					Meta: meta.Metadata{
						Name:        "bound",
						Reference:   "bound",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "bound",
								Reference:   "bound.bound",
								Annotations: map[string]string{},
								Parent:      "",
								UUID:        "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeAppNode",
							},
						},
						Apps: map[string]*meta.App{
							"bound2": {
								Meta: meta.Metadata{
									Name:        "bound2",
									Reference:   "bound.bound2",
									Annotations: map[string]string{},
									Parent:      "bound",
									UUID:        "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "bound2",
											Reference:   "bound.bound2",
											Annotations: map[string]string{},
											Parent:      "bound",
											UUID:        "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeAppNode",
										},
									},
									Apps: map[string]*meta.App{
										"bound3": {
											Meta: meta.Metadata{
												Name:        "bound3",
												Reference:   "bound.bound2.bound3",
												Annotations: map[string]string{},
												Parent:      "bound.bound2",
												UUID:        "",
											},
											Spec: meta.AppSpec{
												Node: meta.Node{
													Meta: meta.Metadata{
														Name:        "bound3",
														Reference:   "bound.bound2.bound3",
														Annotations: map[string]string{},
														Parent:      "bound.bound2",
														UUID:        "",
													},
													Spec: meta.NodeSpec{
														Image: "imageNodeAppNode",
													},
												},
												Apps:     map[string]*meta.App{},
												Channels: map[string]*meta.Channel{},
												Types:    map[string]*meta.Type{},
												Boundary: meta.AppBoundary{
													Input:  []string{"alias1"},
													Output: []string{"alias2"},
												},
											},
										},
									},
									Channels: map[string]*meta.Channel{},
									Types:    map[string]*meta.Type{},
									Boundary: meta.AppBoundary{
										Input:  []string{"alias1"},
										Output: []string{"alias2"},
									},
								},
							},
							"boundNP": {
								Meta: meta.Metadata{
									Name:        "boundNP",
									Reference:   "invalid.path",
									Annotations: map[string]string{},
									Parent:      "invalid.path",
									UUID:        "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "boundNP",
											Reference:   "invalid.path",
											Annotations: map[string]string{},
											Parent:      "bound",
											UUID:        "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeAppNode",
										},
									},
									Apps: map[string]*meta.App{
										"boundNP2": {
											Meta: meta.Metadata{
												Name:        "boundNP2",
												Reference:   "bound.boundNP.boundNP2",
												Annotations: map[string]string{},
												Parent:      "bound.boundNP",
												UUID:        "",
											},
											Spec: meta.AppSpec{
												Node: meta.Node{
													Meta: meta.Metadata{
														Name:        "boundNP2",
														Reference:   "bound.boundNP.boundNP2",
														Annotations: map[string]string{},
														Parent:      "bound.boundNP",
														UUID:        "",
													},
													Spec: meta.NodeSpec{
														Image: "imageNodeAppNode",
													},
												},
												Apps:     map[string]*meta.App{},
												Channels: map[string]*meta.Channel{},
												Types:    map[string]*meta.Type{},
												Boundary: meta.AppBoundary{
													Input:  []string{"alias1"},
													Output: []string{"alias2"},
												},
											},
										},
									},
									Channels: map[string]*meta.Channel{},
									Types:    map[string]*meta.Type{},
									Boundary: meta.AppBoundary{
										Input:  []string{"alias1"},
										Output: []string{"alias2"},
									},
								},
							},
							"bound4": {
								Meta: meta.Metadata{
									Name:        "bound4",
									Reference:   "bound.bound4",
									Annotations: map[string]string{},
									Parent:      "bound",
									UUID:        "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "bound4",
											Reference:   "bound.bound4",
											Annotations: map[string]string{},
											Parent:      "bound",
											UUID:        "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeAppNode",
										},
									},
									Apps:     map[string]*meta.App{},
									Channels: map[string]*meta.Channel{},
									Types:    map[string]*meta.Type{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1"},
										Output: []string{"alias3"},
									},
								},
							},
							"bound5": {
								Meta: meta.Metadata{
									Name:        "bound5",
									Reference:   "bound.bound5",
									Annotations: map[string]string{},
									Parent:      "bound",
									UUID:        "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "bound5",
											Reference:   "bound.bound5",
											Annotations: map[string]string{},
											Parent:      "bound",
											UUID:        "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeAppNode",
										},
									},
									Apps:     map[string]*meta.App{},
									Channels: map[string]*meta.Channel{},
									Types:    map[string]*meta.Type{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1"},
										Output: []string{"alias4"},
									},
								},
							},
							"bound6": {
								Meta: meta.Metadata{
									Name:        "bound6",
									Reference:   "bound.bound6",
									Annotations: map[string]string{},
									Parent:      "bound",
									UUID:        "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "bound6",
											Reference:   "bound.bound6",
											Annotations: map[string]string{},
											Parent:      "bound",
											UUID:        "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeAppNode",
										},
									},
									Apps: map[string]*meta.App{
										"bound7": {
											Meta: meta.Metadata{
												Name:        "bound7",
												Reference:   "bound.bound6",
												Annotations: map[string]string{},
												Parent:      "bound.bound6",
												UUID:        "",
											},
											Spec: meta.AppSpec{
												Node: meta.Node{
													Meta: meta.Metadata{
														Name:        "bound6",
														Reference:   "bound.bound6",
														Annotations: map[string]string{},
														Parent:      "bound",
														UUID:        "",
													},
													Spec: meta.NodeSpec{
														Image: "imageNodeAppNode",
													},
												},
												Apps:     map[string]*meta.App{},
												Channels: map[string]*meta.Channel{},
												Types:    map[string]*meta.Type{},
												Boundary: meta.AppBoundary{
													Input:  []string{"bdch1"},
													Output: []string{"alias3"},
												},
											},
										},
									},
									Channels: map[string]*meta.Channel{},
									Types:    map[string]*meta.Type{},
									Boundary: meta.AppBoundary{
										Input:  []string{"bdch1"},
										Output: []string{"alias3"},
									},
									Aliases: map[string]*meta.Alias{
										"bound8.alias": {
											Target: "notch",
										},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"bdch1": {
								Meta: meta.Metadata{
									Name:   "bdch1",
									Parent: "",
									UUID:   "uuid-bdch1",
								},
								ConnectedApps: []string{},
								Spec:          meta.ChannelSpec{},
							},
							"bdch2": {
								Meta: meta.Metadata{
									Name:   "bdch2",
									Parent: "",
									UUID:   "uuid-bdch2",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						Types: map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
						Aliases: map[string]*meta.Alias{
							"bound2.alias1": {
								Target: "bdch1",
							},
							"bound2.alias2": {
								Target: "bdch2",
							},
							"bound4.alias3": {
								Target: "bdch2",
							},
							"bound6.alias3": {
								Target: "bdch2",
							},
						},
					},
				},
				"connectedApp": {
					Meta: meta.Metadata{
						Name: "connectedApp",
					},
					Spec: meta.AppSpec{
						Apps: map[string]*meta.App{
							"noAliasSon": {
								Meta: meta.Metadata{
									Name:   "noAliasSon",
									Parent: "connectedApp",
								},
								Spec: meta.AppSpec{
									Boundary: meta.AppBoundary{
										Input: utils.StringArray{
											"channel1",
										},
										Output: utils.StringArray{
											"channel2",
										},
									},
								},
							},
							"aliasSon": {
								Meta: meta.Metadata{
									Name:   "aliasSon",
									Parent: "connectedApp",
								},
								Spec: meta.AppSpec{
									Boundary: meta.AppBoundary{
										Input: utils.StringArray{
											"alias1",
										},
										Output: utils.StringArray{
											"alias2S",
										},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"channel1": {
								Meta: meta.Metadata{
									Name: "channel1",
									UUID: "uuid-channel1",
								},
								ConnectedApps: utils.StringArray{
									"noAliasSon",
								},
							},
							"channel2": {
								Meta: meta.Metadata{
									Name: "channel2",
									UUID: "uuid-channel2",
								},
								ConnectedApps: utils.StringArray{
									"noAliasSon",
								},
							},
						},
						Aliases: map[string]*meta.Alias{
							"aliasSon.alias1": {
								Target: "channel1",
							},
							"aliasSon.alias2": {
								Target: "channel2",
							},
						},
					},
				},
				"appForParentInjection": {
					Meta: meta.Metadata{
						Name: "appForParentInjection",
					},
					Spec: meta.AppSpec{
						Apps: make(map[string]*meta.App),
					},
				},
			},
			Channels: map[string]*meta.Channel{
				"ch1": {
					Meta: meta.Metadata{
						Name:   "ch1",
						Parent: "",
						UUID:   "uuid-ch1",
					},
					Spec: meta.ChannelSpec{
						Type: "ct1",
					},
				},
				"ch2": {
					Meta: meta.Metadata{
						Name:   "ch2",
						Parent: "",
						UUID:   "uuid-ch2",
					},
					Spec: meta.ChannelSpec{
						Type: "ct2",
					},
				},
			},
			Types: map[string]*meta.Type{
				"ct1": {
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "root.ct1",
						Annotations: map[string]string{},
						Parent:      "root",
						UUID:        "uuid-ct1",
					},
					Schema: "",
				},
				"ct2": {
					Meta: meta.Metadata{
						Name:        "ct2",
						Reference:   "root.ct2",
						Annotations: map[string]string{},
						Parent:      "root",
						UUID:        "uuid-ct2",
					},
					Schema: "",
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

func TestMemoryManager_Apps(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	tests := []struct {
		name   string
		fields fields
		want   AppMemory
	}{
		{
			name: "creating a AppMemoryManager",
			fields: fields{
				root: getMockApp(),
			},
			want: &AppMemoryManager{
				&treeMemoryManager{
					root: getMockApp(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmm := &treeMemoryManager{
				root: tt.fields.root,
			}
			if got := tmm.Apps(); !metautils.CompareWithUUID(got.(*AppMemoryManager).root, tt.want.(*AppMemoryManager).root) {
				t.Errorf("MemoryManager.Apps() = %v", got)
			}
		})
	}
}

func TestAppMemoryManager_GetApp(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *meta.App
		wantErr bool
	}{
		{
			name: "Getting root app",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				query: "",
			},
			wantErr: false,
			want:    getMockApp(),
		},
		{
			name: "Getting a root's child app",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				query: "app1",
			},
			wantErr: false,
			want:    getMockApp().Spec.Apps["app1"],
		},
		{
			name: "Getting app inside non-root app",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				query: "app2.app3",
			},
			wantErr: false,
			want:    getMockApp().Spec.Apps["app2"].Spec.Apps["app3"],
		},
		{
			name: "Using invalid query",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				query: "app2.app9",
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &treeMemoryManager{
				root: tt.fields.root,
				tree: tt.fields.root,
			}
			amm := mem.Apps()
			got, err := amm.Get(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !metautils.CompareWithoutUUID(got, tt.want) {
				t.Errorf("AppMemoryManager.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppMemoryManager_Create(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	type args struct {
		app         *meta.App
		context     string
		searchQuery string
		brokers     *apimodels.BrokersDI
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		want          *meta.App
		checkFunction func(t *testing.T, tmm *treeMemoryManager)
	}{
		{
			name: "Creating app inside of root",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "",
				searchQuery: "appCr1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "appCr1",
						Reference:   "appCr1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node:     meta.Node{},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			wantErr: false,
			want: &meta.App{
				Meta: meta.Metadata{
					Name:        "appCr1",
					Reference:   "appCr1",
					Annotations: map[string]string{},
					Parent:      "",
					UUID:        "",
				},
				Spec: meta.AppSpec{
					Node:     meta.Node{},
					Apps:     map[string]*meta.App{},
					Channels: map[string]*meta.Channel{},
					Types:    map[string]*meta.Type{},
					Boundary: meta.AppBoundary{
						Input:  []string{},
						Output: []string{},
					},
				},
			},
		},
		{
			name: "Creating app inside of non-root app",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "app2",
				searchQuery: "app2.appCr2-1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "appCr2-1",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node:     meta.Node{},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			wantErr: false,
			want: &meta.App{
				Meta: meta.Metadata{
					Name:        "appCr2-1",
					Reference:   "",
					Annotations: map[string]string{},
					Parent:      "app2",
					UUID:        "",
				},
				Spec: meta.AppSpec{
					Node:     meta.Node{},
					Apps:     map[string]*meta.App{},
					Channels: map[string]*meta.Channel{},
					Types:    map[string]*meta.Type{},
					Boundary: meta.AppBoundary{
						Input:  []string{},
						Output: []string{},
					},
				},
			},
		},
		{
			name: "Creating app with invalid context",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "invalidCtx",
				searchQuery: "invalidCtx.invalidApp",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "invalidApp",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node:     meta.Node{},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Invalid - Creating app inside of app with Node",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "appNode",
				searchQuery: "appNode.appInvalidWithNode",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "appInvalidWithNode",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node:     meta.Node{},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Creating app with conflicting name",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "",
				searchQuery: "app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node:     meta.Node{},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Creating app with existing name but not in the same context",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node:     meta.Node{},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			wantErr: false,
			want: &meta.App{
				Meta: meta.Metadata{
					Name:        "app2",
					Reference:   "",
					Annotations: map[string]string{},
					Parent:      "app2",
					UUID:        "",
				},
				Spec: meta.AppSpec{
					Node:     meta.Node{},
					Apps:     map[string]*meta.App{},
					Channels: map[string]*meta.Channel{},
					Types:    map[string]*meta.Type{},
					Boundary: meta.AppBoundary{
						Input:  []string{},
						Output: []string{},
					},
				},
			},
		},
		{
			name: "Creating app with valid boundary",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "app2.app2",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node:     meta.Node{},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2"},
							Output: []string{"ch2app2"},
						},
					},
				},
			},
			wantErr: false,
			want: &meta.App{
				Meta: meta.Metadata{
					Name:        "app2",
					Reference:   "app2.app2",
					Annotations: map[string]string{},
					Parent:      "app2",
					UUID:        "",
				},
				Spec: meta.AppSpec{
					Node:     meta.Node{},
					Apps:     map[string]*meta.App{},
					Channels: map[string]*meta.Channel{},
					Types:    map[string]*meta.Type{},
					Boundary: meta.AppBoundary{
						Input:  []string{"ch1app2"},
						Output: []string{"ch2app2"},
					},
				},
			},
		},
		{
			name: "Creating app with invalid boundary",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "app2",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node:     meta.Node{},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2invalid"},
							Output: []string{"ch2app2"},
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Creating app with node and other apps in it",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "nodeApp2-2",
								Reference:   "",
								Annotations: map[string]string{},
								Parent:      "app2",
								UUID:        "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeAppTest",
							},
						},
						Apps: map[string]*meta.App{
							"appTest1": {},
						},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Creating app with Node",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "app2",
								Reference:   "",
								Annotations: nil,
								Parent:      "",
								UUID:        "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeAppTest",
							},
						},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			wantErr: false,
			want: &meta.App{
				Meta: meta.Metadata{
					Name:        "app2",
					Reference:   "",
					Annotations: map[string]string{},
					Parent:      "app2",
					UUID:        "",
				},
				Spec: meta.AppSpec{
					Node: meta.Node{
						Meta: meta.Metadata{
							Name:        "app2",
							Reference:   "",
							Annotations: map[string]string{},
							Parent:      "app2",
							UUID:        "",
						},
						Spec: meta.NodeSpec{
							Image: "imageNodeAppTest",
						},
					},
					Apps:     map[string]*meta.App{},
					Channels: map[string]*meta.Channel{},
					Types:    map[string]*meta.Type{},
					Boundary: meta.AppBoundary{
						Input:  []string{},
						Output: []string{},
					},
				},
			},
		},
		{
			name: "It should update the channel's connectedApps list",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context: "app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app7",
						Reference:   "app2.app7",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Apps: map[string]*meta.App{
							"app8": {
								Meta: meta.Metadata{
									Name:        "app8",
									Reference:   "app2.app7.app8",
									Annotations: map[string]string{},
									Parent:      "app2.app7",
									UUID:        "",
								},
								Spec: meta.AppSpec{
									Node:     meta.Node{},
									Apps:     map[string]*meta.App{},
									Channels: map[string]*meta.Channel{},
									Types:    map[string]*meta.Type{},
									Boundary: meta.AppBoundary{
										Input:  []string{"channel1"},
										Output: []string{},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"channel1": {
								Meta: meta.Metadata{
									Name: "channel1",
								},
								ConnectedApps: []string{},
								Spec: meta.ChannelSpec{
									Type: "ct1",
								},
							},
						},
						Types: map[string]*meta.Type{
							"ct1": {
								Meta: meta.Metadata{
									Name: "ct1",
								},
								ConnectedChannels: []string{},
							},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2"},
							Output: []string{"ch2app2"},
						},
					},
				},
			},
			wantErr: false,
			checkFunction: func(t *testing.T, tmm *treeMemoryManager) {
				am := tmm.Channels()
				ch, err := am.Get("app2.app7", "channel1")
				if err != nil {
					t.Errorf("cant get channel channel1")
				}
				if !utils.Includes(ch.ConnectedApps, "app8") {
					fmt.Println(ch.ConnectedApps)
					t.Errorf("connectedApps of channel1 dont have app8")
				}
			},
		},
		{
			name: "Create App with a channel that has a invalid Type",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context: "app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{},
						Channels: map[string]*meta.Channel{
							"channel1": {
								Meta: meta.Metadata{
									Name: "channel1",
								},
								ConnectedApps: []string{},
								Spec: meta.ChannelSpec{
									Type: "invalidTypeName",
								},
							},
						},
						Types: map[string]*meta.Type{
							"ct1": {
								Meta: meta.Metadata{
									Name: "ct1",
								},
								ConnectedChannels: []string{},
							},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2"},
							Output: []string{"ch2app2"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "It should update the Type's connectedChannels list",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context: "app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "app2.app2",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{},
						Channels: map[string]*meta.Channel{
							"channel1": {
								Meta: meta.Metadata{
									Name: "channel1",
								},
								ConnectedApps: []string{},
								Spec: meta.ChannelSpec{
									Type: "ct1",
								},
							},
						},
						Types: map[string]*meta.Type{
							"ct1": {
								Meta: meta.Metadata{
									Name: "ct1",
								},
								ConnectedChannels: []string{},
							},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2"},
							Output: []string{"ch2app2"},
						},
					},
				},
			},
			wantErr: false,
			checkFunction: func(t *testing.T, tmm *treeMemoryManager) {
				am := tmm.Types()
				ct, err := am.Get("app2.app2", "ct1")
				if err != nil {
					t.Errorf("cant get Type ct1")
				}
				if !utils.Includes(ct.ConnectedChannels, "channel1") {
					t.Errorf("connectedChannels of ct1 dont have channel1")
				}
			},
		},
		{
			name: "Invalid name - doesn't create app",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "",
				searchQuery: "appCr1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app%Cr1",
						Reference:   "appCr1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node:     meta.Node{},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Creating app with Node without name",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "",
								Reference:   "",
								Annotations: nil,
								Parent:      "",
								UUID:        "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeAppTest",
							},
						},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
					},
				},
			},
			wantErr: false,
			want: &meta.App{
				Meta: meta.Metadata{
					Name:        "app2",
					Reference:   "",
					Annotations: map[string]string{},
					Parent:      "app2",
					UUID:        "",
				},
				Spec: meta.AppSpec{
					Node: meta.Node{
						Meta: meta.Metadata{
							Name:        "app2",
							Reference:   "",
							Annotations: map[string]string{},
							Parent:      "app2",
							UUID:        "",
						},
						Spec: meta.NodeSpec{
							Image: "imageNodeAppTest",
						},
					},
					Apps:     map[string]*meta.App{},
					Channels: map[string]*meta.Channel{},
					Types:    map[string]*meta.Type{},
					Boundary: meta.AppBoundary{
						Input:  []string{},
						Output: []string{},
					},
				},
			},
		},
		{
			name: "Invalid - App with boundary and channel with same name",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "",
								Reference:   "",
								Annotations: nil,
								Parent:      "",
								UUID:        "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeAppTest",
							},
						},
						Apps: map[string]*meta.App{},
						Channels: map[string]*meta.Channel{
							"output1": {
								Meta: meta.Metadata{
									Name:   "output1",
									Parent: "",
								},
								ConnectedApps: []string{},
								Spec:          meta.ChannelSpec{},
							},
						},
						Types: map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{"input1"},
							Output: []string{"output1"},
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Invalid alias",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "app2",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Aliases: map[string]*meta.Alias{
							"app6.output1": {
								Target: "fakeChannel",
							},
						},
						Apps: map[string]*meta.App{
							"app6": {
								Meta: meta.Metadata{
									Name:        "app6",
									Reference:   "app2.app6",
									Annotations: map[string]string{},
									Parent:      "app2.app2",
									UUID:        "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "app6",
											Reference:   "app6.nodeApp4",
											Annotations: map[string]string{},
											Parent:      "app4",
											UUID:        "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeApp4",
										},
									},
									Boundary: meta.AppBoundary{
										Input:  []string{"output1"},
										Output: []string{"output1"},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"output1": {
								Meta: meta.Metadata{
									Name:   "output1",
									Parent: "",
								},
								ConnectedApps: []string{},
								Spec:          meta.ChannelSpec{},
							},
						},
						Types: map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2"},
							Output: []string{"ch1app2"},
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &treeMemoryManager{
				root: tt.fields.root,
				tree: tt.fields.root,
			}
			am := mem.Apps()
			err := am.Create(tt.args.context, tt.args.app, tt.args.brokers)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				scope, _ := metautils.JoinScopes(tt.args.context, tt.args.app.Meta.Name)
				got, _ := am.Get(scope)
				metautils.RecursiveValidateUUIDS("AppMemoryManager.Create()", got, t)
			}
			if tt.want != nil {
				got, err := am.Get(tt.args.searchQuery)
				if (err != nil) || !metautils.CompareWithoutUUID(got, tt.want) {
					t.Errorf("AppMemoryManager.Get() = %v, want %v", got, tt.want)
				}
			}

			if tt.checkFunction != nil {
				tt.checkFunction(t, mem)
			}
		})
	}
}

func TestAppMemoryManager_Delete(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	type args struct {
		query string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		want          *meta.App
		checkFunction func(t *testing.T, tmm *treeMemoryManager)
	}{
		{
			name: "Deleting leaf app from root",
			fields: fields{
				root: getMockApp(),
			},
			args: args{

				query: "app1.thenewapp",
			},
			wantErr: false,
			want:    nil,
		},
		{
			name: "Deleting leaf app from another app",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				query: "app2.app3",
			},
			wantErr: false,
			want:    nil,
		},
		{
			name: "Deleting app with child apps and channels",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				query: "app2",
			},
			wantErr: false,
			want:    nil,
		},
		{
			name: "Deleting root - invalid deletion",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				query: "",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Deleting with invalid query - invalid deletion",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				query: "invalid.query.to.app",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "It should update the channel's connectedApps",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				query: "app1.thenewapp",
			},
			wantErr: false,
			checkFunction: func(t *testing.T, tmm *treeMemoryManager) {
				am := tmm.Channels()
				ch, err := am.Get("app1", "ch1app1")
				if err != nil {
					t.Errorf("cant get channel ch1app1")
				}
				if utils.Includes(ch.ConnectedApps, "thenewapp") {
					t.Errorf("connectedApps of ch1app1 still have thenewapp")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &treeMemoryManager{
				root: tt.fields.root,
				tree: tt.fields.root,
			}
			am := mem.Apps()
			err := am.Delete(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				got, err := am.Get(tt.args.query)
				if err == nil {
					t.Errorf("AppMemoryManager.Get() = %v, want %v", got, tt.want)
				}
			}
			if tt.checkFunction != nil {
				tt.checkFunction(t, mem)
			}
		})
	}
}

func TestAppMemoryManager_Update(t *testing.T) {
	kafkaConfig := sidecars.KafkaConfig{}
	bmm := brokers.GetBrokerMemory()
	bmm.Create(&kafkaConfig)

	type fields struct {
		root    *meta.App
		updated bool
	}
	type args struct {
		app     *meta.App
		query   string
		brokers *apimodels.BrokersDI
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		want          *meta.App
		checkFunction func(t *testing.T, tmm *treeMemoryManager)
	}{
		{
			name: "invalid- update changing apps' name",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				query: "app1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app1Invalid",
						Reference:   "app1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{},
						Channels: map[string]*meta.Channel{
							"ch1app1": {
								Meta: meta.Metadata{
									Name:   "ch1app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
							"ch2app1": {
								Meta: meta.Metadata{
									Name:   "ch2app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						Types: map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid- updated app has node and child apps",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				query: "app1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app1",
						Reference:   "app1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "nodeApp1",
								Reference:   "app1.nodeApp1",
								Annotations: map[string]string{},
								Parent:      "app1",
								UUID:        "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeApp1",
							},
						},
						Apps: map[string]*meta.App{
							"invalidChildApp": {},
						},
						Channels: map[string]*meta.Channel{
							"ch1app1": {
								Meta: meta.Metadata{
									Name:   "ch1app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
							"ch2app1": {
								Meta: meta.Metadata{
									Name:   "ch2app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						Types: map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "invalid- has structural errors",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				query: "app1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app1Invalid",
						Reference:   "app1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{},
						Channels: map[string]*meta.Channel{
							"ch1app1Invalid": {
								Meta: meta.Metadata{
									Name:   "ch1app1",
									Parent: "app1",
								},
								Spec: meta.ChannelSpec{
									Type: "dsntExist",
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
						Types: map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Valid - updated app doesn't have changes",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				query: "app1",
				app:   getMockApp().Spec.Apps["app1"],
			},
			wantErr: false,
			want:    getMockApp().Spec.Apps["app1"],
		},
		{
			name: "Valid - updated app has changes",
			fields: fields{
				root: getMockApp(),

				updated: true,
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				query: "app1",
				app: &meta.App{
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
									Name:        "appUpdate1",
									Reference:   "app1.appUpdate1",
									Annotations: map[string]string{},
									Parent:      "",
									UUID:        "",
								},
								Spec: meta.AppSpec{},
							},
							"appUpdate2": {
								Meta: meta.Metadata{
									Name:        "appUpdate2",
									Reference:   "app1.appUpdate2",
									Annotations: map[string]string{},
									Parent:      "",
									UUID:        "",
								},
								Spec: meta.AppSpec{},
							},
						},
						Channels: map[string]*meta.Channel{
							"ch1app1": {
								Meta: meta.Metadata{
									Name:   "ch1app1",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
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
							},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
			},
			wantErr: false,
			want: &meta.App{
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
								Name:        "appUpdate1",
								Reference:   "app1.appUpdate1",
								Annotations: map[string]string{},
								Parent:      "app1",
								UUID:        "",
							},
							Spec: meta.AppSpec{},
						},
						"appUpdate2": {
							Meta: meta.Metadata{
								Name:        "appUpdate2",
								Reference:   "app1.appUpdate2",
								Annotations: map[string]string{},
								Parent:      "app1",
								UUID:        "",
							},
							Spec: meta.AppSpec{},
						},
					},
					Channels: map[string]*meta.Channel{
						"ch1app1": {
							Meta: meta.Metadata{
								Name:   "ch1app1",
								Parent: "",
							},
							Spec: meta.ChannelSpec{},
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
							ConnectedChannels: []string{"ch2app1Update"},
						},
					},
					Boundary: meta.AppBoundary{
						Input:  []string{"ch1"},
						Output: []string{"ch2"},
					},
				},
			},
		},
		{
			name: "Valid -  check if connectedApps is updated due to invalid changes in channel structure",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				brokers: &apimodels.BrokersDI{
					Available: []string{"kafka"},
					Default:   "kafka",
				},
				query: "app1",
				app: &meta.App{
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
							"thenewapp": {
								Meta: meta.Metadata{
									Name:        "thenewapp",
									Reference:   "app1.thenewapp",
									Annotations: map[string]string{},
									Parent:      "app1",
									UUID:        "",
								},
								Spec: meta.AppSpec{
									Apps:     map[string]*meta.App{},
									Channels: map[string]*meta.Channel{},
									Types:    map[string]*meta.Type{},
									Boundary: meta.AppBoundary{
										Input:  []string{},
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
									Type: "newType",
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
						Types: map[string]*meta.Type{
							"newType": {
								Meta: meta.Metadata{
									Name:        "newType",
									Reference:   "app1.newType",
									Annotations: map[string]string{},
									Parent:      "app1",
									UUID:        "",
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
			wantErr: false,
			want:    nil,
			checkFunction: func(t *testing.T, tmm *treeMemoryManager) {
				am := tmm.Channels()

				ch, err := am.Get("app1", "ch1app1")
				if err != nil {
					t.Errorf("cant get channel ch1app1")
				}

				if utils.Includes(ch.ConnectedApps, "thenewapp") {
					t.Errorf("connectedApps of ch1app1 still have 'thenewapp'")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &treeMemoryManager{
				root: tt.fields.root,
				tree: tt.fields.root,
			}
			am := mem.Apps()
			err := am.Update(tt.args.query, tt.args.app, tt.args.brokers)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				got, err := am.Get(tt.args.query)
				_, derr := diff.Diff(got, tt.want)
				if derr != nil {
					fmt.Println(derr.Error())
				}

				uuidComp := metautils.CompareWithoutUUID(got, tt.want)
				if (err != nil) || (!uuidComp && !tt.fields.updated) {
					t.Errorf("AppMemoryManager.Get() = %v, want %v", got, tt.want)
				}
			}
			if tt.checkFunction != nil {
				tt.checkFunction(t, mem)
			}
		})
	}
}

func TestAppMemoryManager_ResolveBoundary(t *testing.T) {
	type fields struct {
		MemoryManager *treeMemoryManager
		root          *meta.App
		tree          *meta.App
	}
	type args struct {
		app         *meta.App
		usePermTree bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "Valid - resolve direct boundary",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				app:         getMockApp().Spec.Apps["bound"],
				usePermTree: false,
			},
			want: map[string]string{
				"ch1": "ch1",
				"ch2": "ch2",
			},
			wantErr: false,
		},
		{
			name: "Valid - resolve boundary through alias",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				app:         getMockApp().Spec.Apps["bound"].Spec.Apps["bound2"],
				usePermTree: false,
			},
			want: map[string]string{
				"alias1": "bound.bdch1",
				"alias2": "bound.bdch2",
			},
			wantErr: false,
		},
		{
			name: "Valid - resolve boundary through recursive alias",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				app:         getMockApp().Spec.Apps["bound"].Spec.Apps["bound2"].Spec.Apps["bound3"],
				usePermTree: false,
			},
			want: map[string]string{
				"alias1": "bound.bdch1",
				"alias2": "bound.bdch2",
			},
			wantErr: false,
		},
		{
			name: "Valid - resolve boundary through mixed references",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				app:         getMockApp().Spec.Apps["bound"].Spec.Apps["bound4"],
				usePermTree: false,
			},
			want: map[string]string{
				"ch1":    "ch1",
				"alias3": "bound.bdch2",
			},
			wantErr: false,
		},
		{
			name: "Invalid - app with bad parent",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				app:         getMockApp().Spec.Apps["bound"].Spec.Apps["boundNP"],
				usePermTree: false,
			},
			wantErr: true,
		},
		{
			name: "Invalid - app with bad grandpa",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				app:         getMockApp().Spec.Apps["bound"].Spec.Apps["boundNP"].Spec.Apps["boundNP2"],
				usePermTree: false,
			},
			wantErr: true,
		},
		{
			name: "Use perm tree",
			fields: fields{
				root: getMockApp(),
				tree: getMockApp(),
			},
			args: args{
				app:         getMockApp().Spec.Apps["bound"].Spec.Apps["bound2"],
				usePermTree: true,
			},
			want: map[string]string{
				"alias1": "bound.bdch1",
				"alias2": "bound.bdch2",
			},
			wantErr: false,
		},
		{
			name: "Use perm tree - want error (tree not set)",
			fields: fields{
				root: getMockApp(),
			},
			args: args{
				app:         getMockApp().Spec.Apps["bound"].Spec.Apps["bound2"],
				usePermTree: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &treeMemoryManager{
				root: tt.fields.root,
				tree: tt.fields.tree,
			}
			amm := mem.Apps()
			got, err := amm.ResolveBoundary(tt.args.app, tt.args.usePermTree)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.ResolveBoundary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !metautils.CompareWithoutUUID(got, tt.want) {
				t.Errorf("AppMemoryManager.ResolveBoundary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppMemoryManager_removeFromParentBoundary(t *testing.T) {
	type args struct {
		app    *meta.App
		parent *meta.App
	}
	tests := []struct {
		name string
		args args
		want map[string]*meta.Channel
	}{
		{
			name: "no alias should remove",
			args: args{
				app:    getMockApp().Spec.Apps["connectedApp"].Spec.Apps["noAliasSon"],
				parent: getMockApp().Spec.Apps["connectedApp"],
			},
			want: map[string]*meta.Channel{
				"channel1": {
					Meta: meta.Metadata{
						Name: "channel1",
					},
				},
				"channel2": {
					Meta: meta.Metadata{
						Name: "channel2",
					},
				},
			},
		},
		{
			name: "alias should not remove",
			args: args{
				app:    getMockApp().Spec.Apps["connectedApp"].Spec.Apps["aliasSon"],
				parent: getMockApp().Spec.Apps["connectedApp"],
			},
			want: map[string]*meta.Channel{
				"channel1": {
					Meta: meta.Metadata{
						Name: "channel1",
					},
					ConnectedApps: utils.StringArray{
						"noAliasSon",
					},
				},
				"channel2": {
					Meta: meta.Metadata{
						Name: "channel2",
					},
					ConnectedApps: utils.StringArray{
						"noAliasSon",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &treeMemoryManager{
				root: getMockApp(),
				tree: getMockApp(),
			}
			amm := mem.Apps().(*AppMemoryManager)
			amm.removeFromParentBoundary(tt.args.app, tt.args.parent)
			if !metautils.CompareWithoutUUID(tt.args.parent.Spec.Channels, tt.want) {
				t.Errorf("removeFromParentBoundary() result =\n%#v, want\n%#v", tt.args.parent.Spec.Channels, tt.want)
			}

		})
	}
}
