package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/fake"
	ofake "gitlab.inspr.dev/inspr/core/cmd/insprd/operators/fake"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
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
		App: meta.App{
			Meta: meta.Metadata{
				Name: "mock_app",
			},
		},
		Ctx:   "",
		Valid: true,
	})
	wrongFormatData, _ := json.Marshal([]byte{1})
	return []appAPITest{
		{
			name: "successful_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(errors.New("test_error")), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator()),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().NotFound().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().AlreadyExists().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().InternalServer().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().InvalidName().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().InvalidApp().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannel().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannelType().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().BadRequest().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedAppDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

// appQueryDICases - generates the test cases to be used in functions that
// handle the use the appQueryDI struct of the models package.
// For example, HandleGetAppByRef and HandleDeleteApp use these test cases
func appQueryDICases(funcName string) []appAPITest {
	parsedQueryAppDI, _ := json.Marshal(models.AppQueryDI{
		Ctx:   ".mock_app",
		Valid: true,
	})
	wrongFormatData, _ := json.Marshal([]byte{1})
	return []appAPITest{
		{
			name: "successful_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedQueryAppDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(errors.New("test_error")), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedQueryAppDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "failed_parsing_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator()),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().NotFound().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedQueryAppDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().AlreadyExists().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedQueryAppDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().InternalServer().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedQueryAppDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().InvalidName().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedQueryAppDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().InvalidApp().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedQueryAppDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannel().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedQueryAppDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannelType().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedQueryAppDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			ah:   NewAppHandler(fake.MockMemoryManager(ierrors.NewError().BadRequest().Build()), ofake.NewFakeOperator()),
			send: sendInRequest{body: parsedQueryAppDI},
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
				memManager: fake.MockMemoryManager(nil),
			},
			want: &AppHandler{
				AppMemory: fake.MockMemoryManager(nil).Apps(),
				op:        ofake.NewFakeOperator(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAppHandler(tt.args.memManager, ofake.NewFakeOperator()); !reflect.DeepEqual(got, tt.want) {
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

			tt.ah.CreateApp("", &meta.App{Meta: meta.Metadata{Name: "mock_app"}})

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

			tt.ah.CreateApp("", &meta.App{Meta: meta.Metadata{Name: "mock_app"}})

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
			handlerFunc := tt.ah.HandleDeleteApp()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.ah.CreateApp("", &meta.App{Meta: meta.Metadata{Name: "mock_app"}})

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
