package handler

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

func TestNewChannelTypeHandler(t *testing.T) {
	type args struct {
		memManager memory.Manager
	}
	tests := []struct {
		name string
		args args
		want *ChannelTypeHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChannelTypeHandler(tt.args.memManager); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannelTypeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelTypeHandler_HandleCreateInfo(t *testing.T) {
	tests := []struct {
		name string
		cth  *ChannelTypeHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cth.HandleCreateInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelTypeHandler.HandleCreateInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelTypeHandler_HandleCreateChannelType(t *testing.T) {
	tests := []struct {
		name string
		cth  *ChannelTypeHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cth.HandleCreateChannelType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelTypeHandler.HandleCreateChannelType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelTypeHandler_HandleGetChannelTypeByRef(t *testing.T) {
	tests := []struct {
		name string
		cth  *ChannelTypeHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cth.HandleGetChannelTypeByRef(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelTypeHandler.HandleGetChannelTypeByRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelTypeHandler_HandleUpdateChannelType(t *testing.T) {
	tests := []struct {
		name string
		cth  *ChannelTypeHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cth.HandleUpdateChannelType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelTypeHandler.HandleUpdateChannelType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelTypeHandler_HandleDeleteChannelType(t *testing.T) {
	tests := []struct {
		name string
		cth  *ChannelTypeHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cth.HandleDeleteChannelType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelTypeHandler.HandleDeleteChannelType() = %v, want %v", got, tt.want)
			}
		})
	}
}
