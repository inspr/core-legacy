package tree

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func getMockRootApp() *meta.App {
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
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
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
											Name:        "nodeApp3",
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
										Input:  []string{"ch1"},
										Output: []string{"ch2"},
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
			},
			Channels: map[string]*meta.Channel{
				"ch1": {},
				"ch2": {},
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
					Schema: []byte{},
				},
				"ct2": {
					Meta: meta.Metadata{
						Name:        "ct2",
						Reference:   "root.ct2",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
					},
					Schema: []byte{},
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

func TestTreeMemoryManager_Apps(t *testing.T) {
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
				root: getMockRootApp(),
			},
			want: &AppMemoryManager{
				root: getMockRootApp(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmm := &TreeMemoryManager{
				root: tt.fields.root,
			}
			if got := tmm.Apps(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TreeMemoryManager.Apps() = %v, want %v", got, tt.want)
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
				root:   getMockRootApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				query: "",
			},
			wantErr: false,
			want: &meta.App{
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
								Apps:         map[string]*meta.App{},
								Channels:     map[string]*meta.Channel{},
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
													Name:        "nodeApp3",
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
												Input:  []string{"ch1"},
												Output: []string{"ch2"},
											},
										},
									},
								},
								Channels:     map[string]*meta.Channel{},
								ChannelTypes: map[string]*meta.ChannelType{},
								Boundary: meta.AppBoundary{
									Input:  []string{"ch1"},
									Output: []string{"ch2"},
								},
							},
						},
					},
					Channels: map[string]*meta.Channel{
						"ch1": {},
						"ch2": {},
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
							Schema: []byte{},
						},
						"ct2": {
							Meta: meta.Metadata{
								Name:        "ct2",
								Reference:   "root.ct2",
								Annotations: map[string]string{},
								Parent:      "root",
								SHA256:      "",
							},
							Schema: []byte{},
						},
					},
					Boundary: meta.AppBoundary{
						Input:  []string{},
						Output: []string{},
					},
				},
			},
		},
		{
			name: "Getting a root's child app",
			fields: fields{
				root:   getMockRootApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				query: "app1",
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
					Apps:         map[string]*meta.App{},
					Channels:     map[string]*meta.Channel{},
					ChannelTypes: map[string]*meta.ChannelType{},
					Boundary: meta.AppBoundary{
						Input:  []string{"ch1"},
						Output: []string{"ch2"},
					},
				},
			},
		},
		{
			name: "Getting app inside non-root app",
			fields: fields{
				root:   getMockRootApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				query: "app2.app3",
			},
			wantErr: false,
			want: &meta.App{
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
							Name:        "nodeApp3",
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
						Input:  []string{"ch1"},
						Output: []string{"ch2"},
					},
				},
			},
		},
		{
			name: "Using invalid query",
			fields: fields{
				root:   getMockRootApp(),
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
			setTree(&TreeMockManager{
				root:   tt.fields.root,
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			amm := GetTreeMemory().Apps()
			got, err := amm.GetApp(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.GetApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppMemoryManager.GetApp() = %v, want %v", got, tt.want)
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
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *meta.App
	}{
		{
			name: "Creating app inside of root",
			fields: fields{
				root:   getMockRootApp(),
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
				root:   getMockRootApp(),
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
				root:   getMockRootApp(),
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
			name: "Creating app inside of app with Node",
			fields: fields{
				root:   getMockRootApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				context:     "app1",
				searchQuery: "app1.appInvalidWithNode",
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
				root:   getMockRootApp(),
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
				root:   getMockRootApp(),
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
				root:   getMockRootApp(),
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
				root:   getMockRootApp(),
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
				root:   getMockRootApp(),
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
				root:   getMockRootApp(),
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
			setTree(&TreeMockManager{
				root:   tt.fields.root,
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			am := GetTreeMemory().Apps()
			err := am.CreateApp(tt.args.app, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.CreateApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				got, err := am.GetApp(tt.args.searchQuery)
				if (err != nil) || !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AppMemoryManager.Get() = %v, want %v", got, tt.want)
				}
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
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *meta.App
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTree(&TreeMockManager{
				root:   tt.fields.root,
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			am := GetTreeMemory().Apps()
			err := am.DeleteApp(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.CreateApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				got, err := am.GetApp(tt.args.query)
				if (err != nil) || !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AppMemoryManager.Get() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestAppMemoryManager_UpdateApp(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	type args struct {
		app   *meta.App
		query string
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
			amm := &AppMemoryManager{
				root: tt.fields.root,
			}
			if err := amm.UpdateApp(tt.args.app, tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("AppMemoryManager.UpdateApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func Test_validAppStructure(t *testing.T) {
// 	type args struct {
// 		app       meta.App
// 		parentApp meta.App
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := validAppStructure(tt.args.app, tt.args.parentApp); got != tt.want {
// 				t.Errorf("validAppStructure() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_nodeIsEmpty(t *testing.T) {
	type args struct {
		node meta.Node
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nodeIsEmpty(tt.args.node); got != tt.want {
				t.Errorf("nodeIsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_validBoundaries(t *testing.T) {
// 	type args struct {
// 		bound          meta.AppBoundary
// 		parentChannels map[string]*meta.Channel
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := validBoundaries(tt.args.bound, tt.args.parentChannels); got != tt.want {
// 				t.Errorf("validBoundaries() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_deleteApp(t *testing.T) {
	type args struct {
		app *meta.App
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
			if err := deleteApp(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("deleteApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getParentApp(t *testing.T) {
	type args struct {
		sonQuery string
	}
	tests := []struct {
		name    string
		args    args
		want    *meta.App
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func Test_validUpdateChanges(t *testing.T) {
	type args struct {
		currentApp *meta.App
		newApp     *meta.App
		query      string
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
			if err := validUpdateChanges(tt.args.currentApp, tt.args.newApp, tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("validUpdateChanges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkForChildStructureChanges(t *testing.T) {
	type args struct {
		currentStruct meta.AppSpec
		newStruct     meta.AppSpec
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]Set
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkForChildStructureChanges(tt.args.currentStruct, tt.args.newStruct)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkForChildStructureChanges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkForChildStructureChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_diffError(t *testing.T) {
	type args struct {
		err error
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
			if err := diffError(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("diffError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_invalidChannelChanges(t *testing.T) {
	type args struct {
		changedChannels Set
		newApp          *meta.App
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := invalidChannelChanges(tt.args.changedChannels, tt.args.newApp); got != tt.want {
				t.Errorf("invalidChannelChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
