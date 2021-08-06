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
	"testing"

	"inspr.dev/inspr/pkg/ierrors"
)

func TestClient_Send(t *testing.T) {
	type fields struct {
		c                http.Client
		middleware       Encoder
		decoderGenerator DecoderGenerator
		auth             Authenticator
	}
	type args struct {
		ctx    context.Context
		route  string
		method string
		body   interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantErrReq bool
		want       interface{}
		response   interface{}
	}{
		{
			name: "test post",
			fields: fields{
				c:                http.Client{},
				middleware:       json.Marshal,
				decoderGenerator: JSONDecoderGenerator,
				auth:             nil,
			},
			args: args{
				ctx:    context.Background(),
				route:  "/test",
				method: http.MethodPost,
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
				auth:             nil,
			},
			args: args{
				ctx:    context.Background(),
				route:  "/test",
				method: http.MethodGet,
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
				auth:             nil,
			},
			args: args{
				ctx:    context.Background(),
				route:  "/test",
				method: http.MethodGet,
				body:   "hello",
			},
			wantErr:    true,
			wantErrReq: true,
			want:       "hello",
		},
		{
			name: "middleware error",
			fields: fields{
				c: http.Client{},
				middleware: func(i interface{}) ([]byte, error) {
					return nil, ierrors.New("")
				},
				decoderGenerator: JSONDecoderGenerator,
				auth:             nil,
			},
			args: args{
				ctx:    context.Background(),
				route:  "/test",
				method: http.MethodGet,
				body:   "hello",
			},
			wantErr: true,
			want:    "hello",
		},
		{
			name: "test_auth_token_errorGet",
			fields: fields{
				c:                http.Client{},
				middleware:       json.Marshal,
				decoderGenerator: JSONDecoderGenerator,
				auth:             mockAuth{errGet: errors.New("mock_err")},
			},
			args: args{
				ctx:    context.Background(),
				route:  "/test",
				method: http.MethodPost,
				body:   "hello",
			},
			wantErr: true,

			want: "hello",
		},
		{
			name: "test_auth_token_errorSet",
			fields: fields{
				c:                http.Client{},
				middleware:       json.Marshal,
				decoderGenerator: JSONDecoderGenerator,
				auth:             mockAuth{errSet: errors.New("mock_err")},
			},
			args: args{
				ctx:    context.Background(),
				route:  "/test",
				method: http.MethodPost,
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
					encoder.Encode(ierrors.New(
						"methods are not equal",
					).BadRequest())
					return
				}
				if r.URL.Path != tt.args.route {
					w.WriteHeader(http.StatusNotFound)
					encoder.Encode(ierrors.New(
						"paths are not equal",
					).BadRequest())
					return
				}

				// adds token to response
				w.Header().Add("Authorization", "Bearer mock_token")

				if tt.wantErrReq {
					w.WriteHeader(http.StatusBadRequest)
					encoder.Encode(ierrors.New(
						"wants error",
					).BadRequest())
					return
				}

				var body interface{}

				decoder.Decode(&body)
				if !reflect.DeepEqual(tt.args.body, body) {
					w.WriteHeader(http.StatusBadRequest)
					encoder.Encode(ierrors.New("").BadRequest())
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
				auth:             tt.fields.auth,
			}

			err := c.Send(
				tt.args.ctx,
				tt.args.route,
				tt.args.method,
				tt.args.body,
				&tt.response,
			)
			if (err != nil) != tt.wantErr {
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
			wantMessage: DefaultErr.Error(),
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
						b, _ := json.Marshal(
							ierrors.New("mock_error").Unauthorized())
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantMessage: ierrors.Wrap(
				ierrors.New("mock_error"),
				"status unauthorized",
			).Error(),
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
						b, _ := json.Marshal(
							ierrors.New("mock_error").Forbidden())
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantMessage: ierrors.Wrap(
				ierrors.New("mock_error"),
				"status forbidden",
			).Error(),
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
							Marshal(ierrors.New("this is an error"))
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantMessage: ierrors.New("this is an error").Error(),
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
							ierrors.New("mock_error").Unauthorized(),
						)
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantMessage: ierrors.Wrap(
				ierrors.New("mock_error").Unauthorized(),
				"status unauthorized",
			).Error(),
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
							ierrors.New("mock_error").Forbidden(),
						)
						return ioutil.NopCloser(bytes.NewReader(b))
					}(),
				},
			},
			wantMessage: ierrors.Wrap(
				ierrors.New("mock_error").Unauthorized(),
				"status forbidden",
			).Error(),
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
			if err.Error() != tt.wantMessage {
				t.Errorf(
					"Client.handleResponseErr() error = %v, wantErr %v",
					err.Error(),
					tt.wantMessage,
				)
			}
		})
	}
}
