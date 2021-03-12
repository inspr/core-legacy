package client

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/controller"
	"gitlab.inspr.dev/inspr/core/pkg/controller/mocks"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

func TestNewControllerClient(t *testing.T) {
	type args struct {
		rc *request.Client
	}
	tests := []struct {
		name string
		args args
		want controller.Interface
	}{
		{
			name: "client_creation",
			args: args{
				rc: request.NewJSONClient("mock_url"),
			},
			want: &Client{
				HTTPClient: request.NewJSONClient("mock_url"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewControllerClient(tt.args.rc)

			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf(
					"NewControllerClient() = %v, want %v",
					reflect.TypeOf(got),
					reflect.TypeOf(tt.want),
				)
			}
		})
	}
}

func TestClient_Channels(t *testing.T) {
	check := func(x interface{}) bool {
		// Declare a type object representing ChannelInterface
		channel := reflect.TypeOf((*controller.ChannelInterface)(nil)).Elem()
		// see if implements the channelInterface
		return reflect.PtrTo(reflect.TypeOf(x)).Implements(channel)
	}

	type fields struct {
		rc *request.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   controller.ChannelInterface
	}{
		{
			name: "channels_creation",
			fields: fields{
				rc: request.NewJSONClient("mock_url"),
			},
			want: mocks.NewChannelMock(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewControllerClient(request.NewJSONClient("mock"))

			got := c.Channels()
			if check(got) != check(tt.want) {
				t.Errorf(
					"Client.Channels() = %v, want %v",
					check(got),
					check(tt.want),
				)
			}
		})
	}
}

func TestClient_Apps(t *testing.T) {
	check := func(x interface{}) bool {
		// Declare a type object representing ChannelInterface
		app := reflect.TypeOf((*controller.AppInterface)(nil)).Elem()
		// see if implements the channelInterface
		return reflect.PtrTo(reflect.TypeOf(x)).Implements(app)
	}
	type fields struct {
		rc *request.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   controller.AppInterface
	}{
		{
			name:   "apps_creation",
			fields: fields{rc: request.NewJSONClient("mock")},
			want:   mocks.NewAppMock(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewControllerClient(tt.fields.rc)
			got := c.Apps()

			if check(got) != check(tt.want) {
				t.Errorf(
					"Client.Apps() = %v, want %v",
					check(got),
					check(tt.want),
				)
			}
		})
	}
}

func TestClient_ChannelTypes(t *testing.T) {
	check := func(x interface{}) bool {
		// Declare a type object representing ChannelInterface
		ct := reflect.TypeOf((*controller.ChannelTypeInterface)(nil)).Elem()
		// see if implements the channelInterface
		return reflect.PtrTo(reflect.TypeOf(x)).Implements(ct)
	}
	type fields struct {
		rc *request.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   controller.ChannelTypeInterface
	}{
		{
			name:   "channelType_creation",
			fields: fields{rc: request.NewJSONClient("mock")},
			want:   mocks.NewChannelTypeMock(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewControllerClient(tt.fields.rc)
			got := c.ChannelTypes()

			if check(got) != check(tt.want) {
				t.Errorf(
					"Client.ChannelTypes() = %v, want %v",
					check(got),
					check(tt.want),
				)
			}
		})
	}
}
