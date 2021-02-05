package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

func TestClient_SendRequest(t *testing.T) {
	type fields struct {
		c                http.Client
		middleware       Encoder
		decoderGenerator DecoderGenerator
	}
	type args struct {
		ctx    context.Context
		route  string
		method string
		body   interface{}
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		want     interface{}
		response interface{}
	}{
		{
			name: "test post",
			fields: fields{
				c:                http.Client{},
				middleware:       json.Marshal,
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				ctx:    context.Background(),
				route:  "/test",
				method: "POST",
				body:   "hello",
			},
			wantErr: false,
			want:    "hello",
		},
		{
			name: "test get",
			fields: fields{
				c:                http.Client{},
				middleware:       json.Marshal,
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				ctx:    context.Background(),
				route:  "/test",
				method: "GET",
				body:   "hello",
			},
			wantErr: false,
			want:    "hello",
		},
		{
			name: "test error",
			fields: fields{
				c:                http.Client{},
				middleware:       json.Marshal,
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				ctx:    context.Background(),
				route:  "/test",
				method: "GET",
				body:   "hello",
			},
			wantErr: true,
			want:    "hello",
		},
		{
			name: "middleware error",
			fields: fields{
				c:                http.Client{},
				middleware:       func(i interface{}) ([]byte, error) { return nil, ierrors.NewError().Build() },
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				ctx:    context.Background(),
				route:  "/test",
				method: "GET",
				body:   "hello",
			},
			wantErr: true,
			want:    "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				decoder := json.NewDecoder(r.Body)
				encoder := json.NewEncoder(w)
				if r.Method != tt.args.method {
					w.WriteHeader(400)
					encoder.Encode(ierrors.NewError().BadRequest().Message("methods are not equal").Build())
					return
				}
				if r.URL.Path != tt.args.route {
					w.WriteHeader(404)
					encoder.Encode(ierrors.NewError().BadRequest().Message("paths are not equal").Build())
					return
				}

				if tt.wantErr {
					w.WriteHeader(400)
					encoder.Encode(ierrors.NewError().BadRequest().Message("wants error").Build())
					return
				}

				var body interface{}

				decoder.Decode(&body)
				if !reflect.DeepEqual(tt.args.body, body) {
					w.WriteHeader(400)
					encoder.Encode(ierrors.NewError().BadRequest().Build())
					return
				}
				encoder.Encode(tt.want)
			}

			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()

			c := &Client{
				c:                tt.fields.c,
				baseURL:          s.URL,
				middleware:       tt.fields.middleware,
				decoderGenerator: tt.fields.decoderGenerator,
			}

			if err := c.SendRequest(tt.args.ctx, tt.args.route, tt.args.method, tt.args.body, &tt.response); (err != nil) != tt.wantErr {
				t.Errorf("Client.SendRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.response, tt.want) {
				t.Errorf("Client.SendRequest() response = %v, want %v", tt.response, tt.want)
			}
		})
	}
}

func TestClient_handleResponseErr(t *testing.T) {
	type fields struct {
		c                http.Client
		baseURL          string
		middleware       Encoder
		decoderGenerator DecoderGenerator
	}
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "response with error",
			fields: fields{
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				&http.Response{
					StatusCode: 400,
					Body: func() io.ReadCloser {
						b, _ := json.Marshal(ierrors.NewError().Message("this is an error").Build())
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantErr: true,
		},
		{
			name: "response with other error",
			fields: fields{
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				&http.Response{
					StatusCode: 401,
					Body: func() io.ReadCloser {
						b, _ := json.Marshal(ierrors.NewError().Message("this is an error").Build())
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantErr: true,
		},
		{
			name: "response without error",
			fields: fields{
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				&http.Response{
					StatusCode: 200,
					Body: func() io.ReadCloser {
						b, _ := json.Marshal("hello")
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				c:                tt.fields.c,
				baseURL:          tt.fields.baseURL,
				middleware:       tt.fields.middleware,
				decoderGenerator: tt.fields.decoderGenerator,
			}
			if err := c.handleResponseErr(tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Client.handleResponseErr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
