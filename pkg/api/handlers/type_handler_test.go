package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.uber.org/zap"

	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/memory/fake"
	"inspr.dev/inspr/cmd/insprd/operators"
	ofake "inspr.dev/inspr/cmd/insprd/operators/fake"
	"inspr.dev/inspr/pkg/api/models"
	authmock "inspr.dev/inspr/pkg/auth/mocks"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
)

type TypeAPITest struct {
	name string
	cth  *TypeHandler
	send sendInRequest
	want expectedResponse
}

// TypeDICases - generates the test cases to be used in functions
// that handle the use the TypeDI struct of the models package.
// For example, HandleCreate and HandleUpdate
// use these test cases
func TypeDICases(funcName string) []TypeAPITest {
	parsedCTDI, _ := json.Marshal(models.TypeDI{
		Type:   meta.Type{Meta: meta.Metadata{Name: "mock_Type"}},
		DryRun: false,
	})
	wrongFormatData := []byte{1}
	return []TypeAPITest{
		{
			name: "successful_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(nil, nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(errors.New("test_error"), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(nil, nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().NotFound().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().AlreadyExists().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().InternalServer().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().InvalidName().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().InvalidApp().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().InvalidChannel().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_type_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().InvalidType().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().BadRequest().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

// TypeQueryDICases - generates the test cases to be used in functions
// that handle the use the TypeQueryDI struct of the models package.
// For example, HandleGet and HandleDelete
// use these test cases
func TypeQueryDICases(funcName string) []TypeAPITest {
	parsedCTQDI, _ := json.Marshal(models.TypeQueryDI{
		TypeName: "mock_Type",
		DryRun:   false,
	})
	wrongFormatData := []byte{1}
	return []TypeAPITest{
		{
			name: "successful_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(nil, nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(errors.New("test_error"), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(nil, nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().NotFound().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().AlreadyExists().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().InternalServer().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().InvalidName().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().InvalidApp().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().InvalidChannel().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_type_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().InvalidType().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.GetMockMemoryManager(ierrors.NewError().BadRequest().Build(), nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

func TestNewTypeHandler(t *testing.T) {
	h := NewHandler(
		fake.GetMockMemoryManager(nil, nil),
		ofake.NewFakeOperator(),
		authmock.NewMockAuth(nil),
	)
	type args struct {
		memManager memory.Manager
		op         operators.OperatorInterface
	}
	tests := []struct {
		name string
		args args
		want *TypeHandler
	}{
		{
			name: "success_CreateHandler",
			args: args{
				memManager: fake.GetMockMemoryManager(nil, nil),
				op:         ofake.NewFakeOperator(),
			},
			want: &TypeHandler{
				h,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.NewTypeHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTypeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeHandler_HandleCreate(t *testing.T) {
	logger = zap.New(nil)
	tests := TypeDICases("HandleCreate")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleCreate().HTTPHandlerFunc()
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
				t.Errorf("ChannelHandler.HandleCreate() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestTypeHandler_HandleGet(t *testing.T) {
	logger = zap.New(nil)
	tests := TypeQueryDICases("HandleGet")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleGet().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.Memory.Tree().Types().Create(
				"",
				&meta.Type{Meta: meta.Metadata{Name: "mock_Type"}},
			)

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
				t.Errorf("ChannelHandler.HandleGet() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestTypeHandler_HandleUpdate(t *testing.T) {
	logger = zap.New(nil)
	tests := TypeDICases("HandleUpdate")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleUpdate().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.Memory.Tree().Types().Create("", &meta.Type{Meta: meta.Metadata{Name: "mock_Type"}})

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("ChannelHandler.HandleUpdate() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestTypeHandler_HandleDelete(t *testing.T) {
	logger = zap.New(nil)
	tests := TypeQueryDICases("HandleDelete")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleDelete().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.Memory.Tree().Types().Create("", &meta.Type{Meta: meta.Metadata{Name: "mock_Type"}})

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("ChannelHandler.HandleDelete() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}
