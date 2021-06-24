package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"inspr.dev/inspr/cmd/insprd/memory/fake"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
	ofake "inspr.dev/inspr/cmd/insprd/operators/fake"
	"inspr.dev/inspr/pkg/api/models"
	authmock "inspr.dev/inspr/pkg/auth/mocks"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
)

// sendInRequest is a struct used for all the testing files in this package
// it's contents is a simple { body []byte }
type sendInRequest struct {
	body  []byte
	scope string
}

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
// For example, HandleCreate and HandleUpdate use these test cases
func appDICases(funcName string) []appAPITest {
	parsedAppDI, _ := json.Marshal(models.AppDI{
		App: meta.App{
			Meta: meta.Metadata{
				Name: "mock_app",
			},
		},
	})
	const Scope = ""
	wrongFormatData, _ := json.Marshal([]byte{1})
	return []appAPITest{
		{
			name: "successful_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(errors.New("test_error")), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: wrongFormatData, scope: Scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().NotFound().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().AlreadyExists().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InternalServer().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidName().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidApp().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidChannel().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidType().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().BadRequest().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

// appQueryDICases - generates the test cases to be used in functions that
// handle the use the appQueryDI struct of the models package.
// For example, HandleGet and HandleDelete use these test cases
func appQueryDICases(funcName string) []appAPITest {
	parsedQueryAppDI, _ := json.Marshal(models.AppQueryDI{})
	const scope = "mock_app"
	wrongFormatData, _ := json.Marshal([]byte{1})
	return []appAPITest{
		{
			name: "successful_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(errors.New("test_error")), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "failed_parsing_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: wrongFormatData, scope: scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().NotFound().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().AlreadyExists().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InternalServer().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidName().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidApp().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidChannel().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidType().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			ah:   NewHandler(fake.MockTreeMemory(ierrors.NewError().BadRequest().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

func TestNewAppHandler(t *testing.T) {
	h := NewHandler(
		fake.MockTreeMemory(nil),
		ofake.NewFakeOperator(),
		authmock.NewMockAuth(nil), nil,
	)
	type args struct {
		memManager tree.Manager
	}
	tests := []struct {
		name string
		args args
		want *AppHandler
	}{
		{
			name: "success_TestNewAppHandler",
			args: args{
				memManager: fake.MockTreeMemory(nil),
			},
			want: &AppHandler{h},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.NewAppHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAppHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppHandler_HandleCreate(t *testing.T) {
	tests := appDICases("HandleCreate")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ah.HandleCreate().HTTPHandlerFunc()
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
				t.Errorf("AppHandler.HandleCreate() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestAppHandler_HandleGet(t *testing.T) {
	tests := appQueryDICases("HandleGet")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ah.HandleGet().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.ah.Memory.Apps().Create("", &meta.App{Meta: meta.Metadata{Name: "mock_app"}})

			client := ts.Client()
			req, _ := http.NewRequest(
				http.MethodPost,
				ts.URL,
				bytes.NewBuffer(tt.send.body),
			)
			req.Header.Add("Scope", tt.send.scope)

			res, err := client.Do(req)
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("AppHandler.HandleDelete() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestAppHandler_HandleUpdate(t *testing.T) {
	tests := appDICases("HandleUpdate")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ah.HandleUpdate().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.ah.Memory.Apps().Create("", &meta.App{Meta: meta.Metadata{Name: "mock_app"}})

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("AppHandler.HandleDelete() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestAppHandler_HandleDelete(t *testing.T) {
	tests := appQueryDICases("HandleDelete")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ah.HandleDelete()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.ah.Memory.Apps().Create("", &meta.App{Meta: meta.Metadata{Name: "mock_app"}})

			client := ts.Client()
			req, _ := http.NewRequest(
				http.MethodPost,
				ts.URL,
				bytes.NewBuffer(tt.send.body),
			)
			req.Header.Add("Scope", tt.send.scope)

			res, err := client.Do(req)
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("AppHandler.HandleDelete() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}
