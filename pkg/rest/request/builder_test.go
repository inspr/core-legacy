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
