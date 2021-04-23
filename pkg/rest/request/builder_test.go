package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestClientBuilder_BaseURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		cb   *ClientBuilder
		args args
		want *ClientBuilder
	}{
		{
			name: "base url setting",
			cb: &ClientBuilder{
				c: &Client{},
			},
			args: args{
				url: "test",
			},
			want: &ClientBuilder{
				c: &Client{
					baseURL: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cb.BaseURL(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientBuilder.BaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientBuilder_Encoder(t *testing.T) {
	type args struct {
		encoder Encoder
	}
	tests := []struct {
		name string
		cb   *ClientBuilder
		args args
	}{
		{
			name: "encoder setting",
			cb: &ClientBuilder{
				c: &Client{},
			},
			args: args{
				encoder: json.Marshal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cb.Encoder(tt.args.encoder)
			got, _ := tt.cb.c.encoder("test")
			want, _ := tt.args.encoder("test")
			if !reflect.DeepEqual(got, want) {
				t.Errorf("ClientBuilder.Encoder() = %v, want %v", got, want)
			}
		})
	}
}

func TestClientBuilder_Decoder(t *testing.T) {
	type args struct {
		decoder DecoderGenerator
	}
	tests := []struct {
		name    string
		cb      *ClientBuilder
		args    args
		encoder Encoder
	}{
		{
			name: "decoder setting",
			cb: &ClientBuilder{
				c: &Client{},
			},
			encoder: json.Marshal,
			args: args{
				decoder: JSONDecoderGenerator,
			},
		},
	}
	for _, tt := range tests {

		tt.cb.Decoder(tt.args.decoder)

		var want, got interface{}
		encoded, _ := tt.encoder("test")
		wantDecoder := tt.args.decoder(ioutil.NopCloser(bytes.NewBuffer(encoded)))
		wantDecoder.Decode(&want)

		encoded, _ = tt.encoder("test")
		gotDecoder := tt.cb.c.decoderGenerator(ioutil.NopCloser(bytes.NewBuffer(encoded)))
		gotDecoder.Decode(&got)

		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(got, want) {
				t.Errorf("ClientBuilder.Decoder() = %v, want %v", got, want)
			}
		})
	}
}

func TestClientBuilder_HTTPClient(t *testing.T) {
	type args struct {
		client http.Client
	}
	tests := []struct {
		name string
		cb   *ClientBuilder
		args args
		want *ClientBuilder
	}{
		{
			name: "http client setting",
			cb: &ClientBuilder{
				c: &Client{},
			},
			args: args{
				client: http.Client{
					Timeout: 100,
				},
			},
			want: &ClientBuilder{
				c: &Client{
					c: http.Client{
						Timeout: 100,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cb.HTTPClient(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientBuilder.HTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		want *ClientBuilder
	}{
		{
			name: "basic new client builder",
			want: &ClientBuilder{
				c: &Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientBuilder_Build(t *testing.T) {
	c := &Client{
		c:       *http.DefaultClient,
		baseURL: "this is an url",
	}
	type fields struct {
		c *Client
	}
	tests := []struct {
		name   string
		fields fields
		want   *Client
	}{
		{
			name: "test client creation",
			fields: fields{
				c: c,
			},
			want: c,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := &ClientBuilder{
				c: tt.fields.c,
			}
			if got := cb.Build(); !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("ClientBuilder.Build() = %#v, want %#v", got, tt.want)
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

func TestClientBuilder_Authenticator(t *testing.T) {
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
			cb := &ClientBuilder{
				c: tt.fields.c,
			}

			got := cb.Authenticator(tt.args.au)

			condition := (reflect.TypeOf(got.c.auth) == reflect.TypeOf(nil))

			if condition != tt.wantNil {
				t.Errorf(
					"ClientBuilder.Authenticator() = %v, want nil => %v",
					reflect.TypeOf(got.c.auth),
					tt.wantNil,
				)
			}
		})
	}
}

func TestClientBuilder_Token(t *testing.T) {
	type fields struct {
		c *Client
	}
	type args struct {
		token []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "basic_test_Token",
			fields: fields{
				c: &Client{},
			},
			args: args{
				token: []byte("mock_token"),
			},
			want: map[string]string{
				"Authentication": "Bearer mock_token",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := &ClientBuilder{
				c: tt.fields.c,
			}
			got := cb.Token(tt.args.token)
			if !reflect.DeepEqual(got.c.headers, tt.want) {
				t.Errorf(
					"ClientBuilder.Token() = %v, want %v",
					got.c.headers,
					tt.want,
				)
			}
		})
	}
}

func TestClientBuilder_Header(t *testing.T) {
	type fields struct {
		c *Client
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "header_inserted_values",
			fields: fields{
				c: &Client{},
			},
			args: args{
				key:   "key",
				value: "Bearer token",
			},
			want: map[string]string{
				"key": "Bearer token",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := &ClientBuilder{
				c: tt.fields.c,
			}
			got := cb.Header(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(got.c.headers, tt.want) {
				t.Errorf(
					"ClientBuilder.Header() = %v, want %v",
					got.c.headers,
					tt.want,
				)
			}
		})
	}
}
