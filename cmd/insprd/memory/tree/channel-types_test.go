package tree

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func TestMemoryManager_ChannelTypes(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	tests := []struct {
		name   string
		fields fields
		want   memory.ChannelTypeMemory
	}{
		{
			name: "creating a ChannelTypeMemortMannager",
			fields: fields{
				root: getMockChannelTypes(),
			},
			want: &ChannelTypeMemoryManager{
				root: getMockChannelTypes(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmm := &MemoryManager{
				root: tt.fields.root,
			}
			if got := tmm.ChannelTypes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryManager.ChannelTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelTypeMemoryManager_GetChannelType(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		context string
		ctName  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *meta.ChannelType
		wantErr bool
	}{
		{
			name: "Getting a valid ChannelType on a valid app",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				context: "",
				ctName:  "ct1",
			},
			wantErr: false,
			want: &meta.ChannelType{
				Meta: meta.Metadata{
					Name:        "ct1",
					Reference:   "ct1",
					Annotations: map[string]string{},
					Parent:      "",
					SHA256:      "",
				},
				Schema: "",
			},
		},
		{
			name: "Getting a invalid ChannelType on a valid app",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				context: "",
				ctName:  "ct3",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Getting any ChannelType on a invalid app",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: ierrors.NewError().NotFound().Build(),
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				context: "invalid.context",
				ctName:  "ct42",
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
			ctm := GetTreeMemory().ChannelTypes()
			got, err := ctm.GetChannelType(tt.args.context, tt.args.ctName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChannelTypeMemoryManager.GetChannelType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelTypeMemoryManager.GetChannelType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelTypeMemoryManager_CreateChannelType(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		ct      *meta.ChannelType
		context string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *meta.ChannelType
	}{
		{
			name: "Creating a new ChannelType on a valid app",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.ChannelType{
					Meta: meta.Metadata{
						Name:        "ct3",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Schema: "",
				},
				context: "",
			},
			wantErr: false,
			want: &meta.ChannelType{
				Meta: meta.Metadata{
					Name:        "ct3",
					Reference:   "ct1",
					Annotations: map[string]string{},
					Parent:      "",
					SHA256:      "",
				},
				Schema: "",
			},
		},
		{
			name: "Trying to create an old ChannelType on a valid app",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.ChannelType{
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Schema: "",
				},
				context: "",
			},
			wantErr: true,
			want: &meta.ChannelType{
				Meta: meta.Metadata{
					Name:        "ct1",
					Reference:   "ct1",
					Annotations: map[string]string{},
					Parent:      "",
					SHA256:      "",
				},
				Schema: "",
			},
		},
		{
			name: "Trying to create an ChannelType on a invalid app",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: ierrors.NewError().NotFound().Build(),
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.ChannelType{
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Schema: "",
				},
				context: "invalid.context",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Invalid name - doesn't create channel type",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.ChannelType{
					Meta: meta.Metadata{
						Name:        "-ct3-",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Schema: "",
				},
				context: "",
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
			ctm := GetTreeMemory().ChannelTypes()
			err := ctm.CreateChannelType(tt.args.ct, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChannelTypeMemoryManager.CreateChannelType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				got, err := ctm.GetChannelType(tt.args.context, tt.want.Meta.Name)
				if (err != nil) || !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ChannelTypeMemoryManager.GetChannelType() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestChannelTypeMemoryManager_DeleteChannelType(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		context string
		ctName  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *meta.ChannelType
	}{
		{
			name: "Deleting a valid ChannelType on a valid app",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ctName:  "ct1",
				context: "",
			},
			wantErr: false,
			want:    nil,
		},
		{
			name: "Deleting a invalid ChannelType on a valid app",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ctName:  "ct3",
				context: "",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Deleting any ChannelType on a invalid app",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: ierrors.NewError().NotFound().Build(),
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				context: "invalid.context",
				ctName:  "ct42",
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
			ctm := GetTreeMemory().ChannelTypes()
			if err := ctm.DeleteChannelType(tt.args.context, tt.args.ctName); (err != nil) != tt.wantErr {
				t.Errorf("ChannelTypeMemoryManager.DeleteChannelType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, _ := ctm.GetChannelType(tt.args.context, tt.args.ctName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelTypeMemoryManager.GetChannelType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelTypeMemoryManager_UpdateChannelType(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		ct      *meta.ChannelType
		context string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *meta.ChannelType
	}{
		{
			name: "Updating a valid ChannelType on a valid app",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.ChannelType{
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Schema: string([]byte{0, 1, 0, 1}),
				},
				context: "",
			},
			wantErr: false,
			want: &meta.ChannelType{
				Meta: meta.Metadata{
					Name:        "ct1",
					Reference:   "ct1",
					Annotations: map[string]string{},
					Parent:      "",
					SHA256:      "",
				},
				Schema: string([]byte{0, 1, 0, 1}),
			},
		},
		{
			name: "Updating a invalid ChannelType on a valid app",
			fields: fields{
				root:   getMockChannelTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.ChannelType{
					Meta: meta.Metadata{
						Name:        "ct42",
						Reference:   "ct42",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Schema: "",
				},
				context: "",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Updating any ChannelType on a invalid app",
			fields: fields{
				root:   nil,
				appErr: ierrors.NewError().NotFound().Build(),
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.ChannelType{
					Meta: meta.Metadata{
						Name:        "ct3",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Schema: "",
				},
				context: "invalid.context",
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
			ctm := GetTreeMemory().ChannelTypes()
			if err := ctm.UpdateChannelType(tt.args.ct, tt.args.context); (err != nil) != tt.wantErr {
				t.Errorf("ChannelTypeMemoryManager.UpdateChannelType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				got, err := ctm.GetChannelType(tt.args.context, tt.want.Meta.Name)
				if (err != nil) || !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ChannelTypeMemoryManager.GetChannelType() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func getMockChannelTypes() *meta.App {
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
				"app1": {},
				"app2": {},
			},
			Channels: map[string]*meta.Channel{},
			ChannelTypes: map[string]*meta.ChannelType{
				"ct1": {
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Schema: "",
				},
				"ct2": {
					Meta: meta.Metadata{
						Name:        "ct2",
						Reference:   "ct2",
						Annotations: map[string]string{},
						Parent:      "",
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
