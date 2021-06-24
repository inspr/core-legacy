package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"inspr.dev/inspr/cmd/insprd/memory/fake"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
	"inspr.dev/inspr/cmd/insprd/operators"
	ofake "inspr.dev/inspr/cmd/insprd/operators/fake"
	"inspr.dev/inspr/pkg/api/models"
	authmock "inspr.dev/inspr/pkg/auth/mocks"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
)

type channelAPITest struct {
	name string
	ch   *ChannelHandler
	send sendInRequest
	want expectedResponse
}

// channelDICases - generates the test cases to be used in functions that
// handle the use the channelDI struct of the models package.
// For example, HandleCreate and HandleUpdate use these test cases
func channelDICases(funcName string) []channelAPITest {
	parsedChannelDI, _ := json.Marshal(models.ChannelDI{
		Channel: meta.Channel{
			Meta: meta.Metadata{Name: "mock_channel"},
		},
	})
	wrongFormatData := []byte{1}
	return []channelAPITest{
		{
			name: "successful_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(errors.New("test_error")), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().NotFound().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().AlreadyExists().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InternalServer().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidName().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidApp().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidChannel().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_type_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidType().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().BadRequest().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

// channelQueryDICases - generates the test cases to be used in functions
// that handle the use the channelQueryDI struct of the models package.
// For example, HandleGet and HandleDelete use these test cases
func channelQueryDICases(funcName string) []channelAPITest {
	parsedChannelQueryDI, _ := json.Marshal(models.ChannelQueryDI{
		ChName: "mock_channel",
		DryRun: false,
	})
	wrongFormatData := []byte{1}
	return []channelAPITest{
		{
			name: "successful_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(errors.New("test_error")), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "bad_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(nil), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: wrongFormatData},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "not_found_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().NotFound().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusNotFound},
		},
		{
			name: "already_exists_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().AlreadyExists().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusConflict},
		},
		{
			name: "internal_server_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InternalServer().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusInternalServerError},
		},
		{
			name: "invalid_name_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidName().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_app_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidApp().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_channel_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidChannel().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "invalid_type_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().InvalidType().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusForbidden},
		},
		{
			name: "bad_request_" + funcName,
			ch:   NewHandler(fake.MockTreeMemory(ierrors.NewError().BadRequest().Build()), ofake.NewFakeOperator(), authmock.NewMockAuth(nil), nil).NewChannelHandler(),
			send: sendInRequest{body: parsedChannelQueryDI},
			want: expectedResponse{status: http.StatusBadRequest},
		},
	}
}

func TestNewChannelHandler(t *testing.T) {
	h := NewHandler(
		fake.MockTreeMemory(nil),
		ofake.NewFakeOperator(),
		authmock.NewMockAuth(nil), nil,
	)
	type args struct {
		memManager tree.Manager
		op         operators.OperatorInterface
	}
	tests := []struct {
		name string
		args args
		want *ChannelHandler
	}{
		{
			name: "success_CreateHandler",
			args: args{
				memManager: fake.MockTreeMemory(nil),
				op:         ofake.NewFakeOperator(),
			},
			want: &ChannelHandler{
				h,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.NewChannelHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannelHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelHandler_HandleCreate(t *testing.T) {
	tests := channelDICases("HandleCreate")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ch.HandleCreate().HTTPHandlerFunc()
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

func TestChannelHandler_HandleGet(t *testing.T) {
	tests := channelQueryDICases("HandleGet")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ch.HandleGet().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.ch.Memory.Channels().Create("", &meta.Channel{Meta: meta.Metadata{Name: "mock_channel"}})

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

func TestChannelHandler_HandleUpdate(t *testing.T) {
	tests := channelDICases("HandleUpdate")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ch.HandleUpdate().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.ch.Memory.Channels().Create("", &meta.Channel{Meta: meta.Metadata{Name: "mock_channel"}})
			tt.ch.Operator.Channels().Create(context.Background(), "", &meta.Channel{Meta: meta.Metadata{Name: "mock_channel"}})

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

func TestChannelHandler_HandleDelete(t *testing.T) {
	tests := channelQueryDICases("HandleDelete")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := tt.ch.HandleDelete()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			tt.ch.Memory.Channels().Create("", &meta.Channel{Meta: meta.Metadata{Name: "mock_channel"}})
			tt.ch.Operator.Channels().Create(context.Background(), "", &meta.Channel{Meta: meta.Metadata{Name: "mock_channel"}})

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
