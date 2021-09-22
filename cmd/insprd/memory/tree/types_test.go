package tree

import (
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	metautils "inspr.dev/inspr/pkg/meta/utils"
)

func TestMemoryManager_Types(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	tests := []struct {
		name   string
		fields fields
		want   TypeMemory
	}{
		{
			name: "creating a TypeMemoryManager",
			fields: fields{
				root: getMockTypes(),
			},
			want: &TypeMemoryManager{
				logger,
				&treeMemoryManager{
					root: getMockTypes(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmm := &treeMemoryManager{
				root: tt.fields.root,
			}
			if got := tmm.Types(); !reflect.DeepEqual(got.(*TypeMemoryManager).root, tt.want.(*TypeMemoryManager).root) {
				t.Errorf("MemoryManager.Types() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeMemoryManager_GetType(t *testing.T) {
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
		want    *meta.Type
		wantErr bool
	}{
		{
			name: "Getting a valid Type on a valid app",
			fields: fields{
				root:   getMockTypes(),
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
			want: &meta.Type{
				Meta: meta.Metadata{
					Name:        "ct1",
					Reference:   "ct1",
					Annotations: map[string]string{},
					Parent:      "",
					UUID:        "",
				},
				Schema: "",
			},
		},
		{
			name: "Getting a invalid Type on a valid app",
			fields: fields{
				root:   getMockTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				context: "",
				ctName:  "ct4",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Getting any Type on a invalid app",
			fields: fields{
				root:   getMockTypes(),
				appErr: ierrors.New("").NotFound(),
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
			mem := &MockManager{
				treeMemoryManager: &treeMemoryManager{
					root: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			}
			ctm := mem.Types()
			got, err := ctm.Get(tt.args.context, tt.args.ctName)
			if (err != nil) != tt.wantErr {
				t.Errorf("TypeMemoryManager.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TypeMemoryManager.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeMemoryManager_Create(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		ct      *meta.Type
		context string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *meta.Type
	}{
		{
			name: "Trying to create a new Type with invalid schema",
			fields: fields{
				root:   getMockTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.Type{
					Meta: meta.Metadata{
						Name:        "ct7",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Schema: "",
				},
				context: "",
			},
			wantErr: true,
		},
		{
			name: "Creating a new Type on a valid app",
			fields: fields{
				root:   getMockTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.Type{
					Meta: meta.Metadata{
						Name:        "ct4",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Schema: "{\"type\":\"string\"}",
				},
				context: "",
			},
			wantErr: false,
			want: &meta.Type{
				Meta: meta.Metadata{
					Name:        "ct4",
					Reference:   "ct1",
					Annotations: map[string]string{},
					Parent:      "",
					UUID:        "",
				},
				Schema: "{\"type\":\"string\"}",
			},
		},
		{
			name: "Trying to create an old Type on a valid app",
			fields: fields{
				root:   getMockTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.Type{
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Schema: "{\"type\":\"string\"}",
				},
				context: "",
			},
			wantErr: true,
			want: &meta.Type{
				Meta: meta.Metadata{
					Name:        "ct1",
					Reference:   "ct1",
					Annotations: map[string]string{},
					Parent:      "",
					UUID:        "",
				},
				Schema: "",
			},
		},
		{
			name: "Trying to create an Type on a invalid app",
			fields: fields{
				root:   getMockTypes(),
				appErr: ierrors.New("").NotFound(),
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.Type{
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Schema: "{\"type\":\"string\"}",
				},
				context: "invalid.context",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Invalid name - doesn't create Type",
			fields: fields{
				root:   getMockTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.Type{
					Meta: meta.Metadata{
						Name:        "-ct3-",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Schema: "{\"type\":\"string\"}",
				},
				context: "",
			},
			wantErr: true,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &MockManager{
				treeMemoryManager: &treeMemoryManager{
					root: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			}
			ctm := mem.Types()
			err := ctm.Create(tt.args.context, tt.args.ct)
			if (err != nil) != tt.wantErr {
				t.Errorf("TypeMemoryManager.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				got, err := ctm.Get(tt.args.context, tt.want.Meta.Name)
				if !tt.wantErr {
					if !metautils.ValidateUUID(got.Meta.UUID) {
						t.Errorf("TypeMemoryManager.Create() invalid UUID, uuid=%v", got.Meta.UUID)
					}
				}
				if (err != nil) || !metautils.CompareWithoutUUID(got, tt.want) {
					t.Errorf("TypeMemoryManager.Create() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestTypeMemoryManager_Delete(t *testing.T) {
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
		want    *meta.Type
	}{
		{
			name: "Deleting a valid Type on a valid app",
			fields: fields{
				root:   getMockTypes(),
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
			name: "Deleting a invalid Type on a valid app",
			fields: fields{
				root:   getMockTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ctName:  "ct4",
				context: "",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Deleting any Type on a invalid app",
			fields: fields{
				root:   getMockTypes(),
				appErr: ierrors.New("").NotFound(),
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
		{
			name: "It should not delete the Type because it's been used by a channel",
			fields: fields{
				root:   getMockTypes(),
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
			want: &meta.Type{
				Meta: meta.Metadata{
					Name:        "ct3",
					Reference:   "ct3",
					Annotations: map[string]string{},
					Parent:      "",
					UUID:        "",
				},
				ConnectedChannels: []string{"channel1"},
				Schema:            "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := &MockManager{
				treeMemoryManager: &treeMemoryManager{
					root: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			}
			ctm := mem.Types()
			if err := ctm.Delete(tt.args.context, tt.args.ctName); (err != nil) != tt.wantErr {
				t.Errorf("TypeMemoryManager.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, _ := ctm.Get(tt.args.context, tt.args.ctName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TypeMemoryManager.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeMemoryManager_Update(t *testing.T) {
	type fields struct {
		root   *meta.App
		appErr error
		mockA  bool
		mockC  bool
		mockCT bool
	}
	type args struct {
		ct      *meta.Type
		context string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *meta.Type
	}{
		{
			name: "Updating a valid Type on a valid app",
			fields: fields{
				root:   getMockTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.Type{
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Schema: string([]byte{0, 1, 0, 1}),
				},
				context: "",
			},
			wantErr: false,
			want: &meta.Type{
				Meta: meta.Metadata{
					Name:        "ct1",
					Reference:   "ct1",
					Annotations: map[string]string{},
					Parent:      "",
					UUID:        "",
				},
				Schema: string([]byte{0, 1, 0, 1}),
			},
		},
		{
			name: "Updating a invalid Type on a valid app",
			fields: fields{
				root:   getMockTypes(),
				appErr: nil,
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.Type{
					Meta: meta.Metadata{
						Name:        "ct42",
						Reference:   "ct42",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Schema: "",
				},
				context: "",
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Updating any Type on a invalid app",
			fields: fields{
				root:   nil,
				appErr: ierrors.New("").NotFound(),
				mockA:  true,
				mockC:  true,
				mockCT: false,
			},
			args: args{
				ct: &meta.Type{
					Meta: meta.Metadata{
						Name:        "ct3",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
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
			mem := &MockManager{
				treeMemoryManager: &treeMemoryManager{
					root: tt.fields.root,
				},
				appErr: tt.fields.appErr,
				mockC:  tt.fields.mockC,
				mockA:  tt.fields.mockA,
				mockCT: tt.fields.mockCT,
			}
			ctm := mem.Types()
			if err := ctm.Update(tt.args.context, tt.args.ct); (err != nil) != tt.wantErr {
				t.Errorf("TypeMemoryManager.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				got, err := ctm.Get(tt.args.context, tt.want.Meta.Name)
				if (err != nil) || !metautils.CompareWithUUID(got, tt.want) {
					t.Errorf("TypeMemoryManager.Get() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func getMockTypes() *meta.App {
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
				"app1": {},
				"app2": {},
			},
			Channels: map[string]*meta.Channel{
				"channel1": {
					Meta: meta.Metadata{
						Name:   "channel1",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "ct3",
					},
				},
			},
			Types: map[string]*meta.Type{
				"ct1": {
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Schema: "",
				},
				"ct2": {
					Meta: meta.Metadata{
						Name:        "ct2",
						Reference:   "ct2",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Schema: "",
				},
				"ct3": {
					Meta: meta.Metadata{
						Name:        "ct3",
						Reference:   "ct3",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					ConnectedChannels: []string{"channel1"},
					Schema:            "",
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
