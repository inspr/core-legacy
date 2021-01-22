package handler

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/mocks"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

func TestNewAppHandler(t *testing.T) {
	type args struct {
		memManager memory.Manager
	}
	tests := []struct {
		name string
		args args
		want *AppHandler
	}{
		{
			name: "success - HandleCreateInfo",
			args: args{
				memManager: mocks.MockMemoryManager(nil),
			},
			want: &AppHandler{
				AppMemory: mocks.MockMemoryManager(nil).Apps(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAppHandler(tt.args.memManager); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAppHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppHandler_HandleCreateInfo(t *testing.T) {
	tests := []struct {
		name string
		ah   *AppHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ah.HandleCreateInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppHandler.HandleCreateInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppHandler_HandleCreateApp(t *testing.T) {
	tests := []struct {
		name string
		ah   *AppHandler
		want rest.Handler
	}{
		{
			name: "t",
			ah:   NewAppHandler(mocks.MockMemoryManager(nil)),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ah.HandleCreateApp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppHandler.HandleCreateApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppHandler_HandleGetAppByRef(t *testing.T) {
	tests := []struct {
		name string
		ah   *AppHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ah.HandleGetAppByRef(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppHandler.HandleGetAppByRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppHandler_HandleUpdateApp(t *testing.T) {
	tests := []struct {
		name string
		ah   *AppHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ah.HandleUpdateApp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppHandler.HandleUpdateApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppHandler_HandleDeleteApp(t *testing.T) {
	tests := []struct {
		name string
		ah   *AppHandler
		want rest.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ah.HandleDeleteApp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppHandler.HandleDeleteApp() = %v, want %v", got, tt.want)
			}
		})
	}
}
