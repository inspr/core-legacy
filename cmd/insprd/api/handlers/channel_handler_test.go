package handler

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

func TestNewChannelHandler(t *testing.T) {
	type args struct {
		memManager memory.Manager
	}
	tests := []struct {
		name string
		args args
		want *ChannelHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChannelHandler(tt.args.memManager); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannelHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelHandler_HandleCreateInfo(t *testing.T) {
	tests := []struct {
		name string
		ch   *ChannelHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ch.HandleCreateInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelHandler.HandleCreateInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelHandler_HandleCreateChannel(t *testing.T) {
	tests := []struct {
		name string
		ch   *ChannelHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ch.HandleCreateChannel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelHandler.HandleCreateChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelHandler_HandleGetChannelByRef(t *testing.T) {
	tests := []struct {
		name string
		ch   *ChannelHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ch.HandleGetChannelByRef(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelHandler.HandleGetChannelByRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelHandler_HandleUpdateChannel(t *testing.T) {
	tests := []struct {
		name string
		ch   *ChannelHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ch.HandleUpdateChannel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelHandler.HandleUpdateChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelHandler_HandleDeleteChannel(t *testing.T) {
	tests := []struct {
		name string
		ch   *ChannelHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ch.HandleDeleteChannel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChannelHandler.HandleDeleteChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}
