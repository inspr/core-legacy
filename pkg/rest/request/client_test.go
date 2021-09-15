package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/utils"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		want Client
	}{
		{
			name: "calling_NewClient",
			want: Client{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClient()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_BaseURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		c    Client
		args args
		want Client
	}{
		{
			name: "base url setting",
			c:    Client{},
			args: args{
				url: "test",
			},
			want: Client{
				baseURL: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.BaseURL(tt.args.url)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"Client.BaseURL() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestClient_Encoder(t *testing.T) {
	type args struct {
		encoder Encoder
	}
	tests := []struct {
		name string
		c    *Client
		args args
	}{
		{
			name: "encoder setting",
			c:    NewJSONClient(""),
			args: args{
				encoder: json.Marshal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Encoder(tt.args.encoder)
			got, _ := tt.c.encoder("test")
			want, _ := tt.args.encoder("test")

			if !reflect.DeepEqual(got, want) {
				t.Errorf("Client.Encoder() = %v, want %v", got, want)
			}
		})
	}
}

func TestClient_Decoder(t *testing.T) {
	type args struct {
		decoder DecoderGenerator
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		encoder Encoder
	}{
		{
			name:    "decoder setting",
			c:       NewJSONClient(""),
			encoder: json.Marshal,
			args: args{
				decoder: JSONDecoderGenerator,
			},
		},
	}
	for _, tt := range tests {

		tt.c.Decoder(tt.args.decoder)

		var want, got interface{}
		encoded, _ := tt.encoder("test")
		wantDecoder := tt.args.decoder(
			ioutil.NopCloser(bytes.NewBuffer(encoded)),
		)
		wantDecoder.Decode(&want)

		encoded, _ = tt.encoder("test")
		gotDecoder := tt.c.decoderGenerator(
			ioutil.NopCloser(bytes.NewBuffer(encoded)),
		)
		gotDecoder.Decode(&got)

		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Client.Decoder() = %v, want %v", got, want)
			}
		})
	}
}

func TestClient_HTTPClient(t *testing.T) {
	type args struct {
		client http.Client
	}
	tests := []struct {
		name string
		c    Client
		args args
		want Client
	}{
		{
			name: "http client setting",
			c:    Client{},
			args: args{
				client: http.Client{
					Timeout: 100,
				},
			},
			want: Client{
				c: http.Client{
					Timeout: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.HTTPClient(tt.args.client)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.HTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJSONClient(t *testing.T) {
	type args struct {
		baseURL string
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "basic_newJsonClient_test",
			args: args{baseURL: "mock_url:8080"},
			want: &Client{
				baseURL:          "mock_url:8080",
				encoder:          json.Marshal,
				decoderGenerator: JSONDecoderGenerator,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// checking each component of the client
			got := NewJSONClient(tt.args.baseURL)

			// URL
			if !reflect.DeepEqual(
				got.baseURL,
				tt.want.baseURL,
			) {
				t.Errorf(
					"NewJSONClient() = %v, want %v",
					got.baseURL,
					tt.want.baseURL,
				)
			}

			// encoder
			gotBytes, _ := got.encoder(1)
			wantBytes, _ := tt.want.encoder(1)
			if !reflect.DeepEqual(
				gotBytes,
				wantBytes,
			) {
				t.Errorf(
					"NewJSONClient() = %v, want %v",
					gotBytes,
					wantBytes,
				)
			}

			// decoder
			gotDecoder := got.decoderGenerator(bytes.NewBuffer(gotBytes))
			wantDecoder := tt.want.decoderGenerator(bytes.NewBuffer(wantBytes))
			if !reflect.DeepEqual(
				gotDecoder,
				wantDecoder,
			) {
				t.Errorf(
					"NewJSONClient() = %v, want %v",
					gotDecoder,
					wantDecoder,
				)
			}

		})
	}
}

func TestClient_Authenticator(t *testing.T) {
	type fields struct {
		c *Client
	}
	type args struct {
		au Authenticator
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantNil bool
	}{
		{
			name: "received_auth",
			fields: fields{
				c: &Client{},
			},
			args: args{
				au: mockAuth{
					errGet: nil,
					errSet: nil,
				},
			},
			wantNil: false,
		},
		{
			name: "nil_auth",
			fields: fields{
				c: &Client{},
			},
			args:    args{au: nil},
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := tt.fields.c.Authenticator(tt.args.au)

			condition := (reflect.TypeOf(got.auth) == reflect.TypeOf(nil))

			if condition != tt.wantNil {
				t.Errorf(
					"Client.Authenticator() = %v, want nil => %v",
					reflect.TypeOf(got.auth),
					tt.wantNil,
				)
			}
		})
	}
}

func TestClient_Token(t *testing.T) {
	type fields struct {
		c Client
	}
	type args struct {
		token []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]utils.StringArray
	}{
		{
			name: "basic_test_Token",
			fields: fields{
				c: Client{},
			},
			args: args{
				token: []byte("mock_token"),
			},
			want: map[string]utils.StringArray{
				"Authorization": {"Bearer mock_token"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.c.Token(tt.args.token)
			if !reflect.DeepEqual(got.headers, tt.want) {
				t.Errorf(
					"Client.Token() = %v, want %v",
					got.headers,
					tt.want,
				)
			}
		})
	}
}

func TestClient_Header(t *testing.T) {
	type fields struct {
		c Client
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]utils.StringArray
	}{
		{
			name: "header_inserted_values",
			fields: fields{
				c: Client{},
			},
			args: args{
				key:   "key",
				value: "Bearer token",
			},
			want: map[string]utils.StringArray{
				"key": {"Bearer token"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.c.Header(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(got.headers, tt.want) {
				t.Errorf(
					"Client.Header() = %v, want %v",
					got.headers,
					tt.want,
				)
			}
		})
	}
}

func TestClient_routeToURL(t *testing.T) {
	type args struct {
		route string
	}
	tests := []struct {
		name string
		c    *Client
		args args
		want string
	}{
		{
			name: "basic testing",
			c: &Client{
				baseURL: "http://test",
			},
			args: args{
				route: "/route",
			},
			want: "http://test/route",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.routeToURL(tt.args.route); got != tt.want {
				t.Errorf("Client.routeToURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONDecoderGenerator(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "decoder creation",
			args: args{
				value: "hello",
			},
			want: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, _ := json.Marshal(tt.args.value)
			gotDecoder := JSONDecoderGenerator(bytes.NewBuffer(encoded))
			var got string
			err := gotDecoder.Decode(&got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONDecoderGenerator() = %v, want %v", got, tt.want)
			}
			if err != nil {
				t.Error("error in decoding")
			}
		})
	}
}
