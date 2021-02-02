package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/mocks"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// sendInRequest is a struct used for all the testing files in this package
// it's contents is a simple { body []byte }
type sendInRequest struct{ body []byte }

// expectedResponse is a struct used for all the testing files in this package
// it's contents is a simple { status int }
type expectedResponse struct{ status int }

type appAPITest struct {
	name string
	ah   *AppHandler
	send sendInRequest
	want expectedResponse
}

// appDICases - generates the test cases to be used in functions that handle
// the use the appDI struct of the models package.
// For example, HandleCreateApp and HandleUpdateApp use these test cases
func appDICases(funcName string) []appAPITest {
	parsedAppDI, _ := json.Marshal(models.AppDI{
		App:   meta.App{},
		Ctx:   "",
		Valid: true,
	})
	wrongFormatData, _ := json.Marshal(struct{}{})
	return []appAPITest{
		{
			name: "successful_request_" + funcName,
			ah:   NewAppHandler(mocks.MockMemoryManager(nil)),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ah:   NewAppHandler(mocks.MockMemoryManager(errors.New("test_error"))),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			ah:   NewAppHandler(mocks.MockMemoryManager(nil)),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

// appQueryDICases - generates the test cases to be used in functions that
// handle the use the appQueryDI struct of the models package.
// For example, HandleGetAppByRef and HandleDeleteApp use these test cases
func appQueryDICases(funcName string) []appAPITest {
	parsedAppQueryDI, _ := json.Marshal(models.AppQueryDI{
		Ctx:   "",
		Valid: true,
	})
	wrongFormatData, _ := json.Marshal(struct{}{})
	return []appAPITest{
		{
			name: "successful_request_" + funcName,
			ah:   NewAppHandler(mocks.MockMemoryManager(nil)),
			send: sendInRequest{body: parsedAppQueryDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ah:   NewAppHandler(mocks.MockMemoryManager(errors.New("test_error"))),
			send: sendInRequest{body: parsedAppQueryDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			ah:   NewAppHandler(mocks.MockMemoryManager(nil)),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

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
			name: "success_TestNewAppHandler",
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

func TestAppHandler_HandleCreateApp(t *testing.T) {
	tests := appDICases("HandleCreateApp")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ah.HandleCreateApp().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("AppHandler.HandleCreateApp() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestAppHandler_HandleGetAppByRef(t *testing.T) {
	tests := appQueryDICases("HandleGetAppByRef")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ah.HandleGetAppByRef().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("AppHandler.HandleDeleteApp() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestAppHandler_HandleUpdateApp(t *testing.T) {
	tests := appDICases("HandleUpdateApp")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ah.HandleUpdateApp().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("AppHandler.HandleDeleteApp() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestAppHandler_HandleDeleteApp(t *testing.T) {
	tests := appQueryDICases("HandleDeleteApp")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ah.HandleDeleteApp().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("AppHandler.HandleDeleteApp() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}
