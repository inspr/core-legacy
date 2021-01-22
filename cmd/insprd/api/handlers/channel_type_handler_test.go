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

type channelTypeAPITest struct {
	name string
	cth  *ChannelTypeHandler
	send struct{ reqBody []byte }
	want struct{ status int }
}

func channelTypeDICases(funcName string) []channelTypeAPITest {
	parsedCTDI, _ := json.Marshal(models.ChannelTypeDI{
		ChannelType: meta.ChannelType{},
		Ctx:         "",
	})
	return []channelTypeAPITest{
		{
			name: "successful_request_" + funcName,
			cth:  NewChannelTypeHandler(mocks.MockMemoryManager(nil)),
			send: struct{ reqBody []byte }{reqBody: parsedCTDI},
			want: struct{ status int }{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			cth:  NewChannelTypeHandler(mocks.MockMemoryManager(errors.New("test_error"))),
			send: struct{ reqBody []byte }{reqBody: parsedCTDI},
			want: struct{ status int }{status: http.StatusInternalServerError},
		},
	}
}

func channelTypeQueryDICases(funcName string) []channelTypeAPITest {
	parsedCTQDI, _ := json.Marshal(models.ChannelTypeQueryDI{
		Ctx:    "",
		CtName: "",
	})
	return []channelTypeAPITest{
		{
			name: "successful_request_" + funcName,
			cth:  NewChannelTypeHandler(mocks.MockMemoryManager(nil)),
			send: struct{ reqBody []byte }{reqBody: parsedCTQDI},
			want: struct{ status int }{status: http.StatusOK},
		},
		{
			name: "unsuccessful_request_" + funcName,
			cth:  NewChannelTypeHandler(mocks.MockMemoryManager(errors.New("test_error"))),
			send: struct{ reqBody []byte }{reqBody: parsedCTQDI},
			want: struct{ status int }{status: http.StatusInternalServerError},
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
				memManager: mocks.MockMemoryManager(nil),
			},
			want: &ChannelTypeHandler{
				ChannelTypeMemory: mocks.MockMemoryManager(nil).ChannelTypes(),
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
			handlerFunc := tt.cth.HandleCreateChannelType()
			ts := httptest.NewServer(http.HandlerFunc(handlerFunc))
			defer ts.Close()
			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.reqBody))
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
			handlerFunc := tt.cth.HandleGetChannelTypeByRef()
			ts := httptest.NewServer(http.HandlerFunc(handlerFunc))
			defer ts.Close()
			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.reqBody))
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
			handlerFunc := tt.cth.HandleUpdateChannelType()
			ts := httptest.NewServer(http.HandlerFunc(handlerFunc))
			defer ts.Close()
			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.reqBody))
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
			handlerFunc := tt.cth.HandleDeleteChannelType()
			ts := httptest.NewServer(http.HandlerFunc(handlerFunc))
			defer ts.Close()
			client := ts.Client()
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(tt.send.reqBody))
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
