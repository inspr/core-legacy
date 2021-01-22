package tree

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func TestTreeMemoryManager_ChannelTypes(t *testing.T) {
	type fields struct {
		root *meta.App
	}
	tests := []struct {
		name   string
		fields fields
		want   memory.ChannelTypeMemory
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmm := &TreeMemoryManager{
				root: tt.fields.root,
			}
			if got := tmm.ChannelTypes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TreeMemoryManager.ChannelTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelTypeMemoryManager_CreateChannelType(t *testing.T) {
	type fields struct {
		root *meta.App
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctm := &ChannelTypeMemoryManager{
				root: tt.fields.root,
			}
			if err := ctm.CreateChannelType(tt.args.ct, tt.args.context); (err != nil) != tt.wantErr {
				t.Errorf("ChannelTypeMemoryManager.CreateChannelType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChannelTypeMemoryManager_GetChannelType(t *testing.T) {
	type fields struct {
		root *meta.App
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctm := &ChannelTypeMemoryManager{
				root: tt.fields.root,
			}
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

func TestChannelTypeMemoryManager_DeleteChannelType(t *testing.T) {
	type fields struct {
		root *meta.App
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctm := &ChannelTypeMemoryManager{
				root: tt.fields.root,
			}
			if err := ctm.DeleteChannelType(tt.args.context, tt.args.ctName); (err != nil) != tt.wantErr {
				t.Errorf("ChannelTypeMemoryManager.DeleteChannelType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChannelTypeMemoryManager_UpdateChannelType(t *testing.T) {
	type fields struct {
		root *meta.App
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctm := &ChannelTypeMemoryManager{
				root: tt.fields.root,
			}
			if err := ctm.UpdateChannelType(tt.args.ct, tt.args.context); (err != nil) != tt.wantErr {
				t.Errorf("ChannelTypeMemoryManager.UpdateChannelType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
