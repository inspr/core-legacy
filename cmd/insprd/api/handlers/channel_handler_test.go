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

type channelAPITest struct {
	name string
	ch   *ChannelHandler
	send sendInRequest
	want expectedResponse
}

// channelDICases - generates the test cases to be used in functions that
// handle the use the channelDI struct of the models package.
// For example, HandleCreateChannel and HandleUpdateChannel use these test cases
func channelDICases(funcName string) []channelAPITest {
	parsedChannelDI, _ := json.Marshal(models.ChannelDI{
		Channel: meta.Channel{
			Meta: meta.Metadata{Name: "mock_channel"},
		},
		Ctx:   "",
		Valid: true,
	})
	wrongFormatData := []byte{1}
	return []channelAPITest{
		{
			name: "successful_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(nil)),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(errors.New("test_error"))),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "failed_parsing_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(nil)),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().NotFound().Build())),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().AlreadyExists().Build())),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().InternalServer().Build())),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().InvalidName().Build())),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().InvalidApp().Build())),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannel().Build())),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannelType().Build())),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().BadRequest().Build())),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

// channelQueryDICases - generates the test cases to be used in functions
// that handle the use the channelQueryDI struct of the models package.
// For example, HandleGetChannelByRef and HandleDeleteChannel use these test cases
func channelQueryDICases(funcName string) []channelAPITest {
	parsedChannelQueryDI, _ := json.Marshal(models.ChannelQueryDI{
		Ctx:    "",
		ChName: "mock_channel",
		Valid:  true,
	})
	wrongFormatData := []byte{1}
	return []channelAPITest{
		{
			name: "successful_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(nil)),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(errors.New("test_error"))),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "failed_parsing_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(nil)),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().NotFound().Build())),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().AlreadyExists().Build())),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().InternalServer().Build())),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().InvalidName().Build())),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().InvalidApp().Build())),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannel().Build())),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_type_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().InvalidChannelType().Build())),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			ch:   NewChannelHandler(fake.MockMemoryManager(ierrors.NewError().BadRequest().Build())),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

func TestNewChannelHandler(t *testing.T) {
	type args struct {
		memManager memory.Manager
	}
	tests := []struct {
		name string
		args args
		want *ChannelHandler
	}{
		{
			name: "success_CreateChannelHandler",
			args: args{
				memManager: fake.MockMemoryManager(nil),
			},
			want: &ChannelHandler{
				ChannelMemory: fake.MockMemoryManager(nil).Channels(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChannelHandler(tt.args.memManager); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannelHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelHandler_HandleCreateChannel(t *testing.T) {
	tests := channelDICases("HandleCreateChannel")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ch.HandleCreateChannel().HTTPHandlerFunc()
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
				t.Errorf("ChannelHandler.HandleCreateChannel() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestChannelHandler_HandleGetChannelByRef(t *testing.T) {
	tests := channelQueryDICases("HandleGetChannelByRef")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ch.HandleGetChannelByRef().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.ch.CreateChannel("", &meta.Channel{Meta: meta.Metadata{Name: "mock_channel"}})

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("ChannelHandler.HandleGetChannelByRef() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestChannelHandler_HandleUpdateChannel(t *testing.T) {
	tests := channelDICases("HandleUpdateChannel")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ch.HandleUpdateChannel().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.ch.CreateChannel("", &meta.Channel{Meta: meta.Metadata{Name: "mock_channel"}})

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("ChannelHandler.HandleUpdateChannel() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}

func TestChannelHandler_HandleDeleteChannel(t *testing.T) {
	tests := channelQueryDICases("HandleDeleteChannel")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ch.HandleDeleteChannel().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.ch.CreateChannel("", &meta.Channel{Meta: meta.Metadata{Name: "mock_channel"}})

			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()

			if res.StatusCode != tt.want.status {
				t.Errorf("ChannelHandler.HandleDeleteChannel() = %v, want %v", res.StatusCode, tt.want.status)
			}
		})
	}
}
