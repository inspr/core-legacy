package request

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/inspr/inspr/pkg/ierrors"
)

func TestClient_Send(t *testing.T) {
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
					w.WriteHeader(http.StatusBadRequest)
					encoder.Encode(ierrors.NewError().BadRequest().Message("methods are not equal").Build())
					return
				}
				if r.URL.Path != tt.args.route {
					w.WriteHeader(http.StatusNotFound)
					encoder.Encode(ierrors.NewError().BadRequest().Message("paths are not equal").Build())
					return
				}

				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
					encoder.Encode(ierrors.NewError().BadRequest().Message("wants error").Build())
					return
				}

				var body interface{}

				decoder.Decode(&body)
				if !reflect.DeepEqual(tt.args.body, body) {
					w.WriteHeader(http.StatusBadRequest)
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
				encoder:          tt.fields.middleware,
				decoderGenerator: tt.fields.decoderGenerator,
			}

			if err := c.Send(tt.args.ctx, tt.args.route, tt.args.method, tt.args.body, &tt.response); (err != nil) != tt.wantErr {
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
		name        string
		fields      fields
		args        args
		wantMessage string
	}{
		{
			name: "default error message",
			fields: fields{
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				&http.Response{
					StatusCode: http.StatusBadRequest,
					Body: func() io.ReadCloser {
						b, _ := json.Marshal(nil)
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantMessage: "cannot retrieve error from server",
		},
		{
			name: "default_error_message_unauthorized_code",
			fields: fields{
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				&http.Response{
					StatusCode: http.StatusUnauthorized,
					Body: func() io.ReadCloser {
						b, _ := json.Marshal(nil)
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantMessage: "cannot retrieve error from server",
		},
		{
			name: "default_error_message_forbidden_code",
			fields: fields{
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				&http.Response{
					StatusCode: http.StatusForbidden,
					Body: func() io.ReadCloser {
						b, _ := json.Marshal(nil)
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantMessage: "cannot retrieve error from server",
		},
		{
			name: "response with custom error",
			fields: fields{
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				&http.Response{
					StatusCode: http.StatusBadRequest,
					Body: func() io.ReadCloser {
						b, _ := json.
							Marshal(ierrors.NewError().
								Message("this is an error").
								Build(),
							)
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantMessage: "this is an error",
		},
		{
			name: "response with unauthorized error",
			fields: fields{
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				&http.Response{
					StatusCode: http.StatusUnauthorized,
					Body: func() io.ReadCloser {
						b, _ := json.Marshal(
							ierrors.NewError().
								InnerError(errors.New("mock_error")).
								Unauthorized().
								Build())
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantMessage: "status unauthorized",
		},
		{
			name: "response with forbidden error",
			fields: fields{
				decoderGenerator: JSONDecoderGenerator,
			},
			args: args{
				&http.Response{
					StatusCode: http.StatusForbidden,
					Body: func() io.ReadCloser {
						b, _ := json.Marshal(
							ierrors.NewError().
								InnerError(errors.New("mock_error")).
								Forbidden().
								Build(),
						)
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantMessage: "status forbidden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				c:                tt.fields.c,
				baseURL:          tt.fields.baseURL,
				encoder:          tt.fields.middleware,
				decoderGenerator: tt.fields.decoderGenerator,
			}
			err := c.handleResponseErr(tt.args.resp)
			var got string

			// does it wrap?
			wrapContent := errors.Unwrap(err)

			if wrapContent == nil {
				got = err.Error()
			} else {
				got = wrapContent.Error()
			}

			if strings.TrimSuffix(got, ": ") != tt.wantMessage {
				t.Errorf(
					"Client.handleResponseErr() error = %v, wantErr %v",
					got,
					tt.wantMessage,
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
