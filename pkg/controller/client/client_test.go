package client

import (
	"os"
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/controller"
	"github.com/inspr/inspr/pkg/controller/mocks"
	"github.com/inspr/inspr/pkg/rest/request"
)

func TestNewControllerClient(t *testing.T) {
	type args struct {
		url   string
		token []byte
	}
	tests := []struct {
		name string
		args args
		want controller.Interface
	}{
		{
			name: "client_creation",
			args: args{
				url:   "mock_url",
				token: []byte("token"),
			},
			want: &Client{
				HTTPClient: request.NewJSONClient("mock_url"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ControllerConfig{
				Auth:  nil,
				Scope: "",
				URL:   tt.args.url,
			}
			got := NewControllerClient(config)

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
			config := ControllerConfig{
				Auth:  nil,
				Scope: "",
				URL:   "mock",
			}
			c := NewControllerClient(config)

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
			config := ControllerConfig{
				Auth:  nil,
				Scope: "",
				URL:   "mock",
			}
			c := NewControllerClient(config)
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
			config := ControllerConfig{
				Auth:  nil,
				Scope: "",
				URL:   "mock",
			}
			c := NewControllerClient(config)
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

func TestClient_Alias(t *testing.T) {
	check := func(x interface{}) bool {
		// Declare a type object representing ChannelInterface
		ct := reflect.TypeOf((*controller.AliasInterface)(nil)).Elem()
		// see if implements the channelInterface
		return reflect.PtrTo(reflect.TypeOf(x)).Implements(ct)
	}
	type fields struct {
		rc *request.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   controller.AliasInterface
	}{
		{
			name:   "alias_creation",
			fields: fields{rc: request.NewJSONClient("mock")},
			want:   mocks.NewAliasMock(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ControllerConfig{
				Auth:  nil,
				Scope: "",
				URL:   "mock",
			}
			c := NewControllerClient(config)
			got := c.Alias()

			if check(got) != check(tt.want) {
				t.Errorf(
					"Client.Alias() = %v, want %v",
					check(got),
					check(tt.want),
				)
			}
		})
	}
}

func TestClient_Auth(t *testing.T) {
	check := func(x interface{}) bool {
		// Declare a type object representing AuthorizationInterface
		authI := reflect.TypeOf((*controller.AuthorizationInterface)(nil)).Elem()
		// see if implements the AuthorizationInterface
		return reflect.PtrTo(reflect.TypeOf(x)).Implements(authI)
	}
	type fields struct {
		rc *request.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   controller.AuthorizationInterface
	}{
		{
			name:   "auth_creation",
			fields: fields{rc: request.NewJSONClient("mock")},
			want:   mocks.NewAuthMock(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ControllerConfig{
				Auth:  nil,
				Scope: "",
				URL:   "mock",
			}
			c := NewControllerClient(config)
			got := c.Authorization()

			if check(got) != check(tt.want) {
				t.Errorf(
					"Client.Authorization() = %v, want %v",
					check(got),
					check(tt.want),
				)
			}
		})
	}
}

func TestGetInClusterConfigs(t *testing.T) {
	tests := []struct {
		name     string
		scopeenv string
		tokenenv string
		urlenv   string
		wantErr  bool
	}{
		{
			name:     "error response, invalid INSPR_CONTROLLER_SCOPE",
			tokenenv: "mock_token",
			urlenv:   "mock_url",
			wantErr:  true,
		},
		{
			name:     "error response, invalid INSPR_CONTROLLER_TOKEN",
			scopeenv: "mock_scope",
			urlenv:   "mock_url",
			wantErr:  true,
		},
		{
			name:     "error response, invalid INSPR_INSPRD_ADDRESS",
			scopeenv: "mock_scope",
			tokenenv: "mock_token",
			wantErr:  true,
		},
		{
			name:     "valid response",
			scopeenv: "mock_scope",
			tokenenv: "mock_token",
			urlenv:   "mock_url",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			if tt.scopeenv != "" {
				os.Setenv("INSPR_CONTROLLER_SCOPE", tt.scopeenv)
			}
			if tt.tokenenv != "" {
				os.Setenv("INSPR_CONTROLLER_TOKEN", tt.tokenenv)
			}
			if tt.urlenv != "" {
				os.Setenv("INSPR_INSPRD_ADDRESS", tt.urlenv)
			}
			got, err := GetInClusterConfigs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInClusterConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				want := &ControllerConfig{
					Auth:  authenticator{},
					Scope: tt.scopeenv,
					URL:   tt.urlenv,
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("GetInClusterConfigs() = %v, want %v", got, want)
				}
			}
		})
	}
}
