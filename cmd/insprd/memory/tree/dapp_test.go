package tree

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

func getMockApp() *meta.App {
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
				"appNode": {
					Meta: meta.Metadata{
						Name:        "appNode",
						Reference:   "appNode",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "appNode",
								Reference:   "appNode.appNode",
								Annotations: map[string]string{},
								Parent:      "appNode",
								SHA256:      "",
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
						ChannelTypes: map[string]*meta.ChannelType{},
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
						ChannelTypes: map[string]*meta.ChannelType{},
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
						SHA256:      "",
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
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "app3",
											Reference:   "app3.nodeApp2",
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
							"app4": {
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
											Name:        "app4",
											Reference:   "app4.nodeApp4",
											Annotations: map[string]string{},
											Parent:      "app4",
											SHA256:      "",
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
									ChannelTypes: map[string]*meta.ChannelType{
										"ctapp4": {
											Meta: meta.Metadata{
												Name:        "ctUpdate1",
												Reference:   "app1.ctUpdate1",
												Annotations: map[string]string{},
												Parent:      "app1",
												SHA256:      "",
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
						ChannelTypes: map[string]*meta.ChannelType{},
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
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "bound",
								Reference:   "bound.bound",
								Annotations: map[string]string{},
								Parent:      "",
								SHA256:      "",
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
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "bound2",
											Reference:   "bound.bound2",
											Annotations: map[string]string{},
											Parent:      "bound",
											SHA256:      "",
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
												SHA256:      "",
											},
											Spec: meta.AppSpec{
												Node: meta.Node{
													Meta: meta.Metadata{
														Name:        "bound3",
														Reference:   "bound.bound2.bound3",
														Annotations: map[string]string{},
														Parent:      "bound2",
														SHA256:      "",
													},
													Spec: meta.NodeSpec{
														Image: "imageNodeAppNode",
													},
												},
												Apps:         map[string]*meta.App{},
												Channels:     map[string]*meta.Channel{},
												ChannelTypes: map[string]*meta.ChannelType{},
												Boundary: meta.AppBoundary{
													Input:  []string{"alias1"},
													Output: []string{"alias2"},
												},
											},
										},
									},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
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
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "boundNP",
											Reference:   "invalid.path",
											Annotations: map[string]string{},
											Parent:      "bound",
											SHA256:      "",
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
												SHA256:      "",
											},
											Spec: meta.AppSpec{
												Node: meta.Node{
													Meta: meta.Metadata{
														Name:        "boundNP2",
														Reference:   "bound.boundNP.boundNP2",
														Annotations: map[string]string{},
														Parent:      "bound",
														SHA256:      "",
													},
													Spec: meta.NodeSpec{
														Image: "imageNodeAppNode",
													},
												},
												Apps:         map[string]*meta.App{},
												Channels:     map[string]*meta.Channel{},
												ChannelTypes: map[string]*meta.ChannelType{},
												Boundary: meta.AppBoundary{
													Input:  []string{"alias1"},
													Output: []string{"alias2"},
												},
											},
										},
									},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
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
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "bound4",
											Reference:   "bound.bound4",
											Annotations: map[string]string{},
											Parent:      "bound",
											SHA256:      "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeAppNode",
										},
									},
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
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
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "bound5",
											Reference:   "bound.bound5",
											Annotations: map[string]string{},
											Parent:      "bound",
											SHA256:      "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeAppNode",
										},
									},
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
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
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "bound6",
											Reference:   "bound.bound6",
											Annotations: map[string]string{},
											Parent:      "bound",
											SHA256:      "",
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
												SHA256:      "",
											},
											Spec: meta.AppSpec{
												Node: meta.Node{
													Meta: meta.Metadata{
														Name:        "bound6",
														Reference:   "bound.bound6",
														Annotations: map[string]string{},
														Parent:      "bound",
														SHA256:      "",
													},
													Spec: meta.NodeSpec{
														Image: "imageNodeAppNode",
													},
												},
												Apps:         map[string]*meta.App{},
												Channels:     map[string]*meta.Channel{},
												ChannelTypes: map[string]*meta.ChannelType{},
												Boundary: meta.AppBoundary{
													Input:  []string{"bdch1"},
													Output: []string{"alias3"},
												},
											},
										},
									},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
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
								},
								ConnectedApps: []string{},
								Spec:          meta.ChannelSpec{},
							},
							"bdch2": {
								Meta: meta.Metadata{
									Name:   "bdch2",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{},
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
			},
			Channels: map[string]*meta.Channel{
				"ch1": {
					Meta: meta.Metadata{
						Name:   "ch1",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "ct1",
					},
				},
				"ch2": {
					Meta: meta.Metadata{
						Name:   "ch2",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "ct2",
					},
				},
			},
			ChannelTypes: map[string]*meta.ChannelType{
				"ct1": {
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "root.ct1",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
					},
					Schema: "",
				},
				"ct2": {
					Meta: meta.Metadata{
						Name:        "ct2",
						Reference:   "root.ct2",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
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
		want   memory.AppMemory
	}{
		{
			name: "creating a AppMemoryManager",
			fields: fields{
				root: getMockApp(),
			},
			want: &AppMemoryManager{
				&MemoryManager{
					root: getMockApp(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmm := &MemoryManager{
				root: tt.fields.root,
			}
			if got := tmm.Apps(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryManager.Apps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppMemoryManager_GetApp(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockC  bool
		mockCT bool
		mockA  bool
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
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
			amm := GetTreeMemory().Apps()
			got, err := amm.Get(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppMemoryManager.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppMemoryManager_CreateApp(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		app         *meta.App
		context     string
		searchQuery string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		want          *meta.App
		checkFunction func(t *testing.T)
	}{
		{
			name: "Creating app inside of root",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "",
				searchQuery: "appCr1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "appCr1",
						Reference:   "appCr1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node:         meta.Node{},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
					SHA256:      "",
				},
				Spec: meta.AppSpec{
					Node:         meta.Node{},
					Apps:         map[string]*meta.App{},
					Channels:     map[string]*meta.Channel{},
					ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "app2",
				searchQuery: "app2.appCr2-1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "appCr2-1",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node:         meta.Node{},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
					SHA256:      "",
				},
				Spec: meta.AppSpec{
					Node:         meta.Node{},
					Apps:         map[string]*meta.App{},
					Channels:     map[string]*meta.Channel{},
					ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "invalidCtx",
				searchQuery: "invalidCtx.invalidApp",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "invalidApp",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node:         meta.Node{},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "appNode",
				searchQuery: "appNode.appInvalidWithNode",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "appInvalidWithNode",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node:         meta.Node{},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "",
				searchQuery: "app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node:         meta.Node{},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node:         meta.Node{},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
					SHA256:      "",
				},
				Spec: meta.AppSpec{
					Node:         meta.Node{},
					Apps:         map[string]*meta.App{},
					Channels:     map[string]*meta.Channel{},
					ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node:         meta.Node{},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
					Reference:   "",
					Annotations: map[string]string{},
					Parent:      "app2",
					SHA256:      "",
				},
				Spec: meta.AppSpec{
					Node:         meta.Node{},
					Apps:         map[string]*meta.App{},
					Channels:     map[string]*meta.Channel{},
					ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node:         meta.Node{},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "nodeApp2-2",
								Reference:   "",
								Annotations: map[string]string{},
								Parent:      "app2",
								SHA256:      "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeAppTest",
							},
						},
						Apps: map[string]*meta.App{
							"appTest1": {},
						},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "app2",
								Reference:   "",
								Annotations: nil,
								Parent:      "",
								SHA256:      "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeAppTest",
							},
						},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
					SHA256:      "",
				},
				Spec: meta.AppSpec{
					Node: meta.Node{
						Meta: meta.Metadata{
							Name:        "app2",
							Reference:   "",
							Annotations: map[string]string{},
							Parent:      "app2",
							SHA256:      "",
						},
						Spec: meta.NodeSpec{
							Image: "imageNodeAppTest",
						},
					},
					Apps:         map[string]*meta.App{},
					Channels:     map[string]*meta.Channel{},
					ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context: "app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app7",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"app8": {
								Meta: meta.Metadata{
									Name:        "app8",
									Reference:   "",
									Annotations: map[string]string{},
									Parent:      "",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node:         meta.Node{},
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
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
						ChannelTypes: map[string]*meta.ChannelType{
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
			checkFunction: func(t *testing.T) {
				am := GetTreeMemory().Channels()
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
			name: "Create App with a channel that has a invalid channelType",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context: "app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
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
									Type: "invalidChannelTypeName",
								},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{
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
			name: "It should update the channelType's connectedChannels list",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context: "app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
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
						ChannelTypes: map[string]*meta.ChannelType{
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
			checkFunction: func(t *testing.T) {
				am := GetTreeMemory().ChannelTypes()
				ct, err := am.Get("app2.app2", "ct1")
				if err != nil {
					t.Errorf("cant get channelType ct1")
				}
				if !utils.Includes(ct.ConnectedChannels, "channel1") {
					t.Errorf("connectedChannels of ct1 dont have channel1")
				}
			},
		},
		{
			name: "Invalid name - doesn't create app",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "",
				searchQuery: "appCr1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app%Cr1",
						Reference:   "appCr1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node:         meta.Node{},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "app2",
				searchQuery: "app2.app2",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "",
								Reference:   "",
								Annotations: nil,
								Parent:      "",
								SHA256:      "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeAppTest",
							},
						},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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
					SHA256:      "",
				},
				Spec: meta.AppSpec{
					Node: meta.Node{
						Meta: meta.Metadata{
							Name:        "app2",
							Reference:   "",
							Annotations: map[string]string{},
							Parent:      "app2",
							SHA256:      "",
						},
						Spec: meta.NodeSpec{
							Image: "imageNodeAppTest",
						},
					},
					Apps:         map[string]*meta.App{},
					Channels:     map[string]*meta.Channel{},
					ChannelTypes: map[string]*meta.ChannelType{},
					Boundary: meta.AppBoundary{
						Input:  []string{},
						Output: []string{},
					},
				},
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
			am := GetTreeMemory().Apps()
			err := am.CreateApp(tt.args.context, tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.CreateApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				got, err := am.Get(tt.args.searchQuery)
				if (err != nil) || !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AppMemoryManager.Get() = %v, want %v", got, tt.want)
				}
			}

			if tt.checkFunction != nil {
				tt.checkFunction(t)
			}
		})
	}
}

func TestAppMemoryManager_DeleteApp(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
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
		checkFunction func(t *testing.T)
	}{
		{
			name: "Deleting leaf app from root",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  true,
				mockCT: true,
				mockA:  false,
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  true,
				mockCT: true,
				mockA:  false,
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  true,
				mockCT: true,
				mockA:  false,
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  true,
				mockCT: true,
				mockA:  false,
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  true,
				mockCT: true,
				mockA:  false,
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				query: "app1.thenewapp",
			},
			wantErr: false,
			checkFunction: func(t *testing.T) {
				am := GetTreeMemory().Channels()
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
			am := GetTreeMemory().Apps()
			err := am.DeleteApp(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.DeleteApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				got, err := am.Get(tt.args.query)
				if err == nil {
					t.Errorf("AppMemoryManager.Get() = %v, want %v", got, tt.want)
				}
			}
			if tt.checkFunction != nil {
				tt.checkFunction(t)
			}
		})
	}
}

func TestAppMemoryManager_UpdateApp(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		app   *meta.App
		query string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		want          *meta.App
		checkFunction func(t *testing.T)
	}{
		{
			name: "invalid- update changing apps' name",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				query: "app1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app1Invalid",
						Reference:   "app1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
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
						ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				query: "app1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app1",
						Reference:   "app1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
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
						ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				query: "app1",
				app: &meta.App{
					Meta: meta.Metadata{
						Name:        "app1Invalid",
						Reference:   "app1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
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
						ChannelTypes: map[string]*meta.ChannelType{},
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  true,
				mockCT: true,
				mockA:  false,
			},
			args: args{
				query: "app1",
				app:   getMockApp().Spec.Apps["app1"],
			},
			wantErr: false,
			want:    getMockApp().Spec.Apps["app1"],
		},
		{
			name: "Valid - updated app has changes",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  true,
				mockCT: true,
				mockA:  false,
			},
			args: args{
				query: "app1",
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
							"appUpdate1": {
								Meta: meta.Metadata{
									Name:        "appUpdate1",
									Reference:   "app1.appUpdate1",
									Annotations: map[string]string{},
									Parent:      "",
									SHA256:      "",
								},
								Spec: meta.AppSpec{},
							},
							"appUpdate2": {
								Meta: meta.Metadata{
									Name:        "appUpdate2",
									Reference:   "app1.appUpdate2",
									Annotations: map[string]string{},
									Parent:      "",
									SHA256:      "",
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
						ChannelTypes: map[string]*meta.ChannelType{
							"ctUpdate1": {
								Meta: meta.Metadata{
									Name:        "ctUpdate1",
									Reference:   "app1.ctUpdate1",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
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
					SHA256:      "",
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
								SHA256:      "",
							},
							Spec: meta.AppSpec{},
						},
						"appUpdate2": {
							Meta: meta.Metadata{
								Name:        "appUpdate2",
								Reference:   "app1.appUpdate2",
								Annotations: map[string]string{},
								Parent:      "app1",
								SHA256:      "",
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
					ChannelTypes: map[string]*meta.ChannelType{
						"ctUpdate1": {
							Meta: meta.Metadata{
								Name:        "ctUpdate1",
								Reference:   "app1.ctUpdate1",
								Annotations: map[string]string{},
								Parent:      "app1",
								SHA256:      "",
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				query: "app1",
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
			wantErr: false,
			want:    nil,
			checkFunction: func(t *testing.T) {
				am := GetTreeMemory().Channels()

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
			am := GetTreeMemory().Apps()
			err := am.UpdateApp(tt.args.query, tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.CreateApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				got, err := am.Get(tt.args.query)
				cl, derr := diff.Diff(got, tt.want)
				if derr != nil {
					fmt.Println(derr.Error())
				}
				cl.Print(os.Stdout)
				if (err != nil) || !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AppMemoryManager.Get() = %v, want %v", got, tt.want)
				}
			}
			if tt.checkFunction != nil {
				tt.checkFunction(t)
			}
		})
	}
}

func TestAppMemoryManager_ResolveBoundary(t *testing.T) {
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
		want    map[string]string
		wantErr bool
	}{
		{
			name: "Valid - resolve direct boundary",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				app: getMockApp().Spec.Apps["bound"],
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
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				app: getMockApp().Spec.Apps["bound"].Spec.Apps["bound2"],
			},
			want: map[string]string{
				"alias1": "bdch1",
				"alias2": "bdch2",
			},
			wantErr: false,
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
			want: map[string]string{
				"alias1": "bdch1",
				"alias2": "bdch2",
			},
			wantErr: false,
		},
		{
			name: "Valid - resolve boundary through mixed references",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				app: getMockApp().Spec.Apps["bound"].Spec.Apps["bound4"],
			},
			want: map[string]string{
				"ch1":    "ch1",
				"alias3": "bdch2",
			},
			wantErr: false,
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
			name: "Invalid - app with bad grandpa",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				app: getMockApp().Spec.Apps["bound"].Spec.Apps["boundNP"].Spec.Apps["boundNP2"],
			},
			wantErr: true,
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
			amm := GetTreeMemory().Apps()
			got, err := amm.ResolveBoundary(tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.ResolveBoundary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppMemoryManager.ResolveBoundary() = %v, want %v", got, tt.want)
			}
		})
	}
}
