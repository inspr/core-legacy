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
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	ofake "gitlab.inspr.dev/inspr/core/cmd/insprd/operators/fake"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type channelTypeAPITest struct {
	name string
	cth  *ChannelTypeHandler
	send sendInRequest
	want expectedResponse
}

// channelTypeDICases - generates the test cases to be used in functions
// that handle the use the channelTypeDI struct of the models package.
// For example, HandleCreate and HandleUpdate
// use these test cases
func channelTypeDICases(funcName string) []channelTypeAPITest {
	parsedCTDI, _ := json.Marshal(models.ChannelTypeDI{
		ChannelType: meta.ChannelType{Meta: meta.Metadata{Name: "mock_channelType"}},
		Ctx:         "",
		Valid:       true,
		DryRun:      false,
	})
	wrongFormatData := []byte{1}
	return []channelTypeAPITest{
		{
			name: "successful_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(errors.New("test_error")), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().NotFound().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().AlreadyExists().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InternalServer().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidName().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidApp().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannel().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannelType().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().BadRequest().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

// channelTypeQueryDICases - generates the test cases to be used in functions
// that handle the use the ChannelTypeQueryDI struct of the models package.
// For example, HandleGetChannelTypeByRef and HandleDelete
// use these test cases
func channelTypeQueryDICases(funcName string) []channelTypeAPITest {
	parsedCTQDI, _ := json.Marshal(models.ChannelTypeQueryDI{
		Ctx:    "",
		CtName: "mock_channelType",
		Valid:  true,
		DryRun: false,
	})
	wrongFormatData := []byte{1}
	return []channelTypeAPITest{
		{
			name: "successful_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(errors.New("test_error")), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(nil), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().NotFound().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().AlreadyExists().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InternalServer().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidName().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidApp().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannel().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannelType().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewHandler(fake.MockMemoryManager(ierrors.NewError().BadRequest().Build()), ofake.NewFakeOperator()).NewChannelTypeHandler(),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

func TestNewChannelTypeHandler(t *testing.T) {
	h := NewHandler(
		fake.MockMemoryManager(nil),
		ofake.NewFakeOperator(),
	)
	type args struct {
		memManager memory.Manager
		op         operators.OperatorInterface
	}
	tests := []struct {
		name string
		args args
		want *ChannelTypeHandler
	}{
		{
			name: "success_CreateHandler",
			args: args{
				memManager: fake.MockMemoryManager(nil),
				op:         ofake.NewFakeOperator(),
			},
			want: &ChannelTypeHandler{
				h,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.NewChannelTypeHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannelTypeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelTypeHandler_HandleCreate(t *testing.T) {
	tests := channelTypeDICases("HandleCreate")
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

func TestChannelTypeHandler_HandleGetChannelTypeByRef(t *testing.T) {
	tests := channelTypeQueryDICases("HandleGetChannelTypeByRef")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleGet().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.Memory.ChannelTypes().Create("", &meta.ChannelType{Meta: meta.Metadata{Name: "mock_channelType"}})

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("ChannelHandler.HandleGetChannelTypeByRef() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestChannelTypeHandler_HandleUpdate(t *testing.T) {
	tests := channelTypeDICases("HandleUpdate")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleUpdate().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.Memory.ChannelTypes().Create("", &meta.ChannelType{Meta: meta.Metadata{Name: "mock_channelType"}})

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

func TestChannelTypeHandler_HandleDelete(t *testing.T) {
	tests := channelTypeQueryDICases("HandleDelete")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleDelete().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.Memory.ChannelTypes().Create("", &meta.ChannelType{Meta: meta.Metadata{Name: "mock_channelType"}})

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
