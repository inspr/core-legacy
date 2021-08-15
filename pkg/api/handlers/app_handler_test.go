package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/memory/fake"
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
			ah: NewHandler(
				fake.GetMockMemoryManager(nil, nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((errors.New("test_error")), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager(nil, nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: wrongFormatData, scope: Scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((ierrors.New("").NotFound()), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager(
					(ierrors.New("").AlreadyExists()),
					nil,
				),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager(
					(ierrors.New("").InternalServer()),
					nil,
				),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((ierrors.New("").InvalidName()), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((ierrors.New("").InvalidApp()), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager(
					(ierrors.New("").InvalidChannel()),
					nil,
				),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((ierrors.New("").InvalidType()), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedAppDI, scope: Scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((ierrors.New("").BadRequest()), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
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
			ah: NewHandler(
				fake.GetMockMemoryManager(nil, nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((errors.New("test_error")), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "failed_parsing_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager(nil, nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: wrongFormatData, scope: scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((ierrors.New("").NotFound()), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager(
					(ierrors.New("").AlreadyExists()),
					nil,
				),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager(
					(ierrors.New("").InternalServer()),
					nil,
				),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((ierrors.New("").InvalidName()), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((ierrors.New("").InvalidApp()), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager(
					(ierrors.New("").InvalidChannel()),
					nil,
				),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((ierrors.New("").InvalidType()), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			ah: NewHandler(
				fake.GetMockMemoryManager((ierrors.New("").BadRequest()), nil),
				ofake.NewFakeOperator(),
				authmock.NewMockAuth(nil),
			).NewAppHandler(),
			send: sendInRequest{body: parsedQueryAppDI, scope: scope},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

func TestNewAppHandler(t *testing.T) {
	h := NewHandler(
		fake.GetMockMemoryManager(nil, nil),
		ofake.NewFakeOperator(),
		authmock.NewMockAuth(nil),
	)
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
				memManager: fake.GetMockMemoryManager(nil, nil),
			},
			want: &AppHandler{h, logger},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.NewAppHandler(); !reflect.DeepEqual(
				got.Handler,
				tt.want.Handler,
			) {
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
			res, err := client.Post(
				ts.URL,
				"application/json",
				bytes.NewBuffer(tt.send.body),
			)
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf(
					"AppHandler.HandleCreate() = %v, want %v",
					res.StatusCode,
					tt.want.status,
				)
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

			tt.ah.Memory.Tree().
				Apps().
				Create("", &meta.App{Meta: meta.Metadata{Name: "mock_app"}}, &models.BrokersDI{})

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
				t.Errorf(
					"AppHandler.HandleDelete() = %v, want %v",
					res.StatusCode,
					tt.want.status,
				)
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

			tt.ah.Memory.Tree().
				Apps().
				Create("", &meta.App{Meta: meta.Metadata{Name: "mock_app"}}, &models.BrokersDI{})

			client := ts.Client()
			res, err := client.Post(
				ts.URL,
				"application/json",
				bytes.NewBuffer(tt.send.body),
			)
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf(
					"AppHandler.HandleDelete() = %v, want %v",
					res.StatusCode,
					tt.want.status,
				)
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

			tt.ah.Memory.Tree().
				Apps().
				Create("", &meta.App{Meta: meta.Metadata{Name: "mock_app"}}, &models.BrokersDI{})

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
				t.Errorf(
					"AppHandler.HandleDelete() = %v, want %v",
					res.StatusCode,
					tt.want.status,
				)
			}
		})
	}
}
