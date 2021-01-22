package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/mocks"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
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
	parsedAppDI, _ := json.Marshal(models.AppDI{
		App: meta.App{},
		Ctx: "",
	})
	tests := []struct {
		name string
		ah   *AppHandler
		send struct {
			reqBody []byte
		}
		want struct {
			status int
			err    error
		}
	}{
		{
			name: "successful_request",
			ah:   NewAppHandler(mocks.MockMemoryManager(nil)),
			send: struct{ reqBody []byte }{
				reqBody: parsedAppDI,
			},
			want: struct {
				status int
				err    error
			}{
				status: http.StatusOK,
				err:    nil,
			},
		},
		{
			name: "unsuccessful_request",
			ah:   NewAppHandler(mocks.MockMemoryManager(nil)),
			send: struct{ reqBody []byte }{
				reqBody: parsedAppDI,
			},
			want: struct {
				status int
				err    error
			}{
				status: http.StatusOK,
				err:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ah.HandleCreateApp()
			ts := httptest.NewServer(http.HandlerFunc((handlerFunc)))
			defer ts.Close()
			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.reqBody))
			if err != tt.want.err {
				t.Errorf("AppHandler.HandleDeleteApp() = %v, want %v", err, tt.want.err)
			}
			defer res.Body.Close()
			if res.StatusCode != tt.want.status {
				t.Errorf("AppHandler.HandleDeleteApp() = %v, want %v", http.StatusOK, tt.want.status)
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
		{
			name: "working handler",
			ah:   NewAppHandler(mocks.MockMemoryManager(nil)),
			want: func(w http.ResponseWriter, r *http.Request) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
