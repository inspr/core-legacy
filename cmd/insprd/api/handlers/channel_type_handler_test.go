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
// For example, HandleCreateChannelType and HandleUpdateChannelType
// use these test cases
func channelTypeDICases(funcName string) []channelTypeAPITest {
	parsedCTDI, _ := json.Marshal(models.ChannelTypeDI{
		ChannelType: meta.ChannelType{Meta: meta.Metadata{Name: "mock_channelType"}},
		Ctx:         "",
		Valid:       true,
	})
	wrongFormatData := []byte{1}
	return []channelTypeAPITest{
		{
			name: "successful_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(nil)),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(errors.New("test_error"))),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "failed_parsing_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(nil)),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().NotFound().Build())),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().AlreadyExists().Build())),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().InternalServer().Build())),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().InvalidName().Build())),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().InvalidApp().Build())),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannel().Build())),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannelType().Build())),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().BadRequest().Build())),
			send: sendInRequest{body: parsedCTDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

// channelTypeQueryDICases - generates the test cases to be used in functions
// that handle the use the ChannelTypeQueryDI struct of the models package.
// For example, HandleGetChannelTypeByRef and HandleDeleteChannelType
// use these test cases
func channelTypeQueryDICases(funcName string) []channelTypeAPITest {
	parsedCTQDI, _ := json.Marshal(models.ChannelTypeQueryDI{
		Ctx:    "",
		CtName: "mock_channelType",
		Valid:  true,
	})
	wrongFormatData := []byte{1}
	return []channelTypeAPITest{
		{
			name: "successful_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(nil)),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(errors.New("test_error"))),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "failed_parsing_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(nil)),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().NotFound().Build())),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().AlreadyExists().Build())),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().InternalServer().Build())),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().InvalidName().Build())),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().InvalidApp().Build())),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannel().Build())),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannelType().Build())),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			cth:  NewChannelTypeHandler(fake.MockMemoryManager(ierrors.NewError().BadRequest().Build())),
			send: sendInRequest{body: parsedCTQDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

func TestNewChannelTypeHandler(t *testing.T) {
	type args struct {
		memManager memory.Manager
	}
	tests := []struct {
		name string
		args args
		want *ChannelTypeHandler
	}{
		{
			name: "success_CreateChannelHandler",
			args: args{
				memManager: fake.MockMemoryManager(nil),
			},
			want: &ChannelTypeHandler{
				ChannelTypeMemory: fake.MockMemoryManager(nil).ChannelTypes(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChannelTypeHandler(tt.args.memManager); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannelTypeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelTypeHandler_HandleCreateChannelType(t *testing.T) {
	tests := channelTypeDICases("HandleCreateChannelType")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleCreateChannelType().HTTPHandlerFunc()
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
				t.Errorf("ChannelHandler.HandleCreateChannelType() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestChannelTypeHandler_HandleGetChannelTypeByRef(t *testing.T) {
	tests := channelTypeQueryDICases("HandleGetChannelTypeByRef")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleGetChannelTypeByRef().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.CreateChannelType("", &meta.ChannelType{Meta: meta.Metadata{Name: "mock_channelType"}})

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

func TestChannelTypeHandler_HandleUpdateChannelType(t *testing.T) {
	tests := channelTypeDICases("HandleUpdateChannelType")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleUpdateChannelType().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.CreateChannelType("", &meta.ChannelType{Meta: meta.Metadata{Name: "mock_channelType"}})

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("ChannelHandler.HandleUpdateChannelType() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestChannelTypeHandler_HandleDeleteChannelType(t *testing.T) {
	tests := channelTypeQueryDICases("HandleDeleteChannelType")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.cth.HandleDeleteChannelType().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.cth.CreateChannelType("", &meta.ChannelType{Meta: meta.Metadata{Name: "mock_channelType"}})

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("ChannelHandler.HandleDeleteChannelType() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}
