package tree

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
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
											Name:        "nodeApp4",
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
				root: getMockApp(),
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
			name: "Creating app inside of app with Node",
			fields: fields{
				root:   getMockApp(),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTree(&MockManager{
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
				query: "app1",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTree(&MockManager{
				root:   tt.fields.root,
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
				got, err := am.GetApp(tt.args.query)
				if err == nil {
					t.Errorf("AppMemoryManager.Get() = %v, want %v", got, tt.want)
				}
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
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *meta.App
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
							"appUpdate1": {},
							"appUpdate2": {},
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
						"appUpdate1": {},
						"appUpdate2": {},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTree(&MockManager{
				root:   tt.fields.root,
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			am := GetTreeMemory().Apps()
			err := am.UpdateApp(tt.args.app, tt.args.query)
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
			name: "invalidapp name - equal to another existing app",
			args: args{
				app: meta.App{
					Meta: meta.Metadata{
						Name:        "app3",
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
				parentApp: *getMockApp().Spec.Apps["app1"],
			},
			want: "parent has Node;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validAppStructure(tt.args.app, tt.args.parentApp); got != tt.want {
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
				bound: meta.AppBoundary{
					Input:  []string{"ch1app1"},
					Output: []string{"ch2app1"},
				},
				parentChannels: getMockApp().Spec.Apps["app1"].Spec.Channels,
			},
			want: "",
		},
		{
			name: "invalidboundary - parent without channels",
			args: args{
				bound: meta.AppBoundary{
					Input:  []string{"ch1app2"},
					Output: []string{"ch2app2"},
				},
				parentChannels: getMockApp().Spec.Apps["app2"].Spec.Apps["app3"].Spec.Channels,
			},
			want: "parent doesn't have Channels;",
		},
		{
			name: "invalid input and output boundary",
			args: args{
				bound: meta.AppBoundary{
					Input:  []string{"ch1app1"},
					Output: []string{"ch2app1"},
				},
				parentChannels: getMockApp().Spec.Apps["app2"].Spec.Channels,
			},
			want: "invalid input boundary;invalid output boundary;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validBoundaries(tt.args.bound, tt.args.parentChannels); got != tt.want {
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
				root:   tt.fields.root,
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

func Test_validUpdateChanges(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockC  bool
		mockCT bool
		mockA  bool
	}
	type args struct {
		currentApp *meta.App
		newApp     *meta.App
		query      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "App without changes",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentApp: getMockApp().Spec.Apps["app2"].Spec.Apps["app4"],
				newApp:     getMockApp().Spec.Apps["app2"].Spec.Apps["app4"],
				query:      "app2.app4",
			},
			wantErr: false,
		},
		{
			name: "Valid boundary changes",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentApp: getMockApp().Spec.Apps["app2"].Spec.Apps["app4"],
				query:      "app2.app4",
				newApp: &meta.App{
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
						Boundary: meta.AppBoundary{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalidboundary changes",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentApp: getMockApp().Spec.Apps["app2"].Spec.Apps["app4"],
				query:      "app2.app4",
				newApp: &meta.App{
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
							Input:  []string{"InvalidInput"},
							Output: []string{"ch2app2"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid structural changes - app with node",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentApp: getMockApp().Spec.Apps["app2"].Spec.Apps["app4"],
				query:      "app2.app4",
				newApp: &meta.App{
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
								Parent:      "app4",
								SHA256:      "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeAppChanged4",
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
							"newch2app4": {
								Meta: meta.Metadata{
									Name:   "newch2app4",
									Parent: "app4",
								},
								Spec: meta.ChannelSpec{
									Type: "ctapp4",
								},
							},
							"newch3app4": {
								Meta: meta.Metadata{
									Name:   "newch3app4",
									Parent: "app4",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{
							"ctapp4": {
								Meta: meta.Metadata{
									Name:        "ctUpdate1",
									Reference:   "app4.ctUpdate1",
									Annotations: map[string]string{},
									Parent:      "app4",
									SHA256:      "",
								},
							},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2"},
							Output: []string{"ch1app2"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid structural changes - app with subapps",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentApp: getMockApp().Spec.Apps["app2"].Spec.Apps["app4"],
				query:      "app2.app4",
				newApp: &meta.App{
					Meta: meta.Metadata{
						Name:        "app4",
						Reference:   "app2.app4",
						Annotations: map[string]string{},
						Parent:      "app2",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"newApp4-1": {
								Meta: meta.Metadata{
									Name:        "newApp4-1",
									Reference:   "app4.newApp4-1",
									Annotations: map[string]string{},
									Parent:      "app4",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node:         meta.Node{},
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary:     meta.AppBoundary{},
								},
							},
						},
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
							"newch2app4": {
								Meta: meta.Metadata{
									Name:   "newch2app4",
									Parent: "app4",
								},
								Spec: meta.ChannelSpec{
									Type: "ctapp4",
								},
							},
							"newch3app4": {
								Meta: meta.Metadata{
									Name:   "newch3app4",
									Parent: "app4",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{
							"ctapp4": {
								Meta: meta.Metadata{
									Name:        "ctUpdate1",
									Reference:   "app4.ctUpdate1",
									Annotations: map[string]string{},
									Parent:      "app4",
									SHA256:      "",
								},
							},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2"},
							Output: []string{"ch1app2"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalidchannel changes",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentApp: getMockApp().Spec.Apps["app2"].Spec.Apps["app4"],
				query:      "app2.app4",
				newApp: &meta.App{
					Meta: meta.Metadata{
						Name:        "app4",
						Reference:   "app2.app4",
						Annotations: map[string]string{},
						Parent:      "app2",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"newApp4-1": {},
							"newApp4-2": {},
						},
						Channels: map[string]*meta.Channel{
							"newch3app4": {
								Meta: meta.Metadata{
									Name:   "newch3app4",
									Parent: "app4",
								},
								Spec: meta.ChannelSpec{
									Type: "invalidChannelType",
								},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary:     meta.AppBoundary{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalidapp changes",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentApp: getMockApp().Spec.Apps["app2"].Spec.Apps["app4"],
				query:      "app2.app4",
				newApp: &meta.App{
					Meta: meta.Metadata{
						Name:        "app4",
						Reference:   "app2.app4",
						Annotations: map[string]string{},
						Parent:      "app2",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"newApp4-1": {
								Meta: meta.Metadata{
									Name:        "newApp4-1",
									Reference:   "app4.newApp4-1",
									Annotations: map[string]string{},
									Parent:      "app4",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node:         meta.Node{},
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary: meta.AppBoundary{
										Input: []string{"invalidInput"},
									},
								},
							},
						},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary:     meta.AppBoundary{},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTree(&MockManager{
				root:   tt.fields.root,
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
			if err := validUpdateChanges(*tt.args.currentApp, *tt.args.newApp, tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("validUpdateChanges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkForChildStructureChanges(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockC  bool
		mockCT bool
		mockA  bool
	}
	type args struct {
		currentStruct meta.AppSpec
		newStruct     meta.AppSpec
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]Set
		wantErr bool
	}{
		{
			name: "No changes",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentStruct: getMockApp().Spec.Apps["app1"].Spec,
				newStruct:     getMockApp().Spec.Apps["app1"].Spec,
			},
			wantErr: false,
			want: map[string]Set{
				"app":     {},
				"channel": {},
				"ctype":   {},
			},
		},
		{
			name: "Channel changes",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentStruct: getMockApp().Spec.Apps["app1"].Spec,
				newStruct: meta.AppSpec{
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
					Apps: map[string]*meta.App{},
					Channels: map[string]*meta.Channel{
						"ch1app1": {
							Meta: meta.Metadata{
								Name:   "ch1app1",
								Parent: "",
							},
							Spec: meta.ChannelSpec{},
						},
						"ch2app1Changed": {
							Meta: meta.Metadata{
								Name:   "ch2app1Changed",
								Parent: "",
							},
							Spec: meta.ChannelSpec{},
						},
						"newChannelApp1": {
							Meta: meta.Metadata{
								Name:   "newChannelApp1",
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
			wantErr: false,
			want: map[string]Set{
				"app": {},
				"channel": {
					"ch2app1Changed": true,
					"newChannelApp1": true,
				},
				"ctype": {},
			},
		},
		{
			name: "ChannelType changes",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentStruct: getMockApp().Spec.Apps["app1"].Spec,
				newStruct: meta.AppSpec{
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
						Input:  []string{"ch1"},
						Output: []string{"ch2"},
					},
				},
			},
			wantErr: false,
			want: map[string]Set{
				"app":     {},
				"channel": {},
				"ctype": {
					"newChannelType": true,
				},
			},
		},
		{
			name: "Child app changes",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentStruct: getMockApp().Spec.Apps["app1"].Spec,
				newStruct: meta.AppSpec{
					Apps: map[string]*meta.App{
						"newChangedApp": {
							Meta: meta.Metadata{
								Name:        "newChangedApp",
								Reference:   "app1.newChangedApp",
								Annotations: map[string]string{},
								Parent:      "app1",
								SHA256:      "",
							},
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
			wantErr: false,
			want: map[string]Set{
				"app": {
					"newChangedApp": true,
				},
				"channel": {},
				"ctype":   {},
			},
		},
		{
			name: "All structures change",
			fields: fields{
				root:   getMockApp(),
				appErr: nil,
				mockC:  false,
				mockCT: false,
				mockA:  false,
			},
			args: args{
				currentStruct: getMockApp().Spec.Apps["app1"].Spec,
				newStruct: meta.AppSpec{
					Apps: map[string]*meta.App{
						"newChangedApp": {
							Meta: meta.Metadata{
								Name:        "newChangedApp",
								Reference:   "app1.newChangedApp",
								Annotations: map[string]string{},
								Parent:      "app1",
								SHA256:      "",
							},
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
						"ch2app1Changed": {
							Meta: meta.Metadata{
								Name:   "ch2app1Changed",
								Parent: "",
							},
							Spec: meta.ChannelSpec{},
						},
						"newChannelApp1": {
							Meta: meta.Metadata{
								Name:   "newChannelApp1",
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
						Input:  []string{"ch1"},
						Output: []string{"ch2"},
					},
				},
			},
			wantErr: false,
			want: map[string]Set{
				"app": {
					"newChangedApp": true,
				},
				"channel": {
					"ch2app1Changed": true,
					"newChannelApp1": true,
				},
				"ctype": {
					"newChannelType": true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTree(&MockManager{
				root:   tt.fields.root,
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			})
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
		{
			name: "Returns ierror",
			args: args{
				err: ierrors.NewError().Build(),
			},
			wantErr: true,
		},
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
		{
			name: "Valid channel change",
			args: args{
				changedChannels: Set{
					"ch1app4": true,
				},
				newApp: getMockApp().Spec.Apps["app2"].Spec.Apps["app4"],
			},
			want: false,
		},
		{
			name: "App is channel's parent but doesn't have channel type",
			args: args{
				changedChannels: Set{
					"ch1app4": true,
				},
				newApp: &meta.App{
					Meta: meta.Metadata{
						Name:        "app4",
						Reference:   "app4",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
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
						},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1app2"},
							Output: []string{"ch2app2"},
						},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := invalidChannelChanges(tt.args.changedChannels, tt.args.newApp); got != tt.want {
				t.Errorf("invalidChannelChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
