package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/memory/fake"
	"github.com/inspr/inspr/cmd/insprd/operators"
	ofake "github.com/inspr/inspr/cmd/insprd/operators/fake"
	"github.com/inspr/inspr/pkg/api/models"
	authmock "github.com/inspr/inspr/pkg/auth/mocks"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
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
		Type:   meta.Type{Meta: meta.Metadata{Name: "mock_channelType"}},
		DryRun: false,
	})
	wrongFormatData := []byte{1}
	return []TypeAPITest{
		{
			name: "successful_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(errors.New("test_error")), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().NotFound().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().AlreadyExists().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InternalServer().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidName().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidApp().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannel().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_type_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidType().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().BadRequest().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
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
		CtName: "mock_channelType",
		DryRun: false,
	})
	wrongFormatData := []byte{1}
	return []TypeAPITest{
		{
			name: "successful_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(errors.New("test_error")), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().NotFound().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().AlreadyExists().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InternalServer().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidName().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidApp().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannel().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_type_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidType().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().BadRequest().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil)).NewTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

func TestNewTypeHandler(t *testing.T) {
	h := NewHandler(
		fake.MockMemoryManager(nil),
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
				memManager: fake.MockMemoryManager(nil),
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
	tests := TypeQueryDICases("HandleGet")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleGet().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.Memory.Types().Create("", &meta.Type{Meta: meta.Metadata{Name: "mock_Type"}})

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
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
	tests := TypeDICases("HandleUpdate")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleUpdate().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.Memory.Types().Create("", &meta.Type{Meta: meta.Metadata{Name: "mock_Type"}})

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
	tests := TypeQueryDICases("HandleDelete")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleDelete().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.Memory.Types().Create("", &meta.Type{Meta: meta.Metadata{Name: "mock_Type"}})

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
