package rest

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/inspr/inspr/pkg/auth"
	authMock "github.com/inspr/inspr/pkg/auth/mocks"
	"github.com/inspr/inspr/pkg/ierrors"
)

// emptyHandler - to be used in testing of the package
var emptyHandler Handler = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestHandler_JSON(t *testing.T) {
	// manipulation - calls .JSON() of the handler in the parameter
	var manipulation = func(h Handler) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			h.JSON()(w, r)
		}
	}
	tests := []struct {
		routeName     string
		expected      string
		customHandler Handler
	}{
		{
			routeName:     "/success",
			expected:      "application/json",
			customHandler: manipulation(emptyHandler),
		},
		{
			routeName:     "/fail",
			expected:      "",
			customHandler: emptyHandler,
		},
	}

	for _, tt := range tests {
		// sets up the test server
		req, err := http.NewRequest("GET", tt.routeName, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		successHandler := http.HandlerFunc(tt.customHandler)
		successHandler.ServeHTTP(rr, req)

		if ct := rr.Header().Get("Content-Type"); ct != tt.expected {
			t.Errorf("Handler.JSON() = %v, want %v", ct, tt.expected)
		}
	}
}

func TestHandler_Validate(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	type args struct {
		auth          auth.Auth
		Authorization string
		httpMethod    string
		reqURLSuffix  string
		scope         string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantCode int
	}{
		{
			name: "no_auth_header",
			args: args{
				auth:       authMock.NewMockAuth(nil),
				httpMethod: http.MethodPost,
			},
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "invalid_token",
			args: args{
				auth: authMock.NewMockAuth(
					&ierrors.InsprError{Code: ierrors.InvalidToken},
				),
				Authorization: "Bearer mock_token",
				httpMethod:    http.MethodPost,
			},
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "unknown_error",
			args: args{
				auth: authMock.NewMockAuth(
					errors.New("mock_error"),
				),
				Authorization: "Bearer mock_token",
				httpMethod:    http.MethodPost,
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "scope_error",
			args: args{
				auth:          authMock.NewMockAuth(nil),
				Authorization: "Bearer mock_token",
				scope:         "no_valid_scope",
				httpMethod:    http.MethodPost,
			},
			wantErr:  true,
			wantCode: http.StatusForbidden,
		},
		{
			name: "scope_error_doesnt_have_prefix",
			args: args{
				auth:          authMock.NewMockAuth(nil),
				Authorization: "Bearer mock_token",
				scope:         "wrongScope.scope_1",
				httpMethod:    http.MethodPost,
			},
			wantErr:  true,
			wantCode: http.StatusForbidden,
		},
		{
			name: "no_permission_for_operation",
			args: args{
				auth:          authMock.NewMockAuth(nil),
				Authorization: "Bearer mock_token",
				scope:         "scope_1",
				httpMethod:    http.MethodGet,
				reqURLSuffix:  "channels",
			},
			wantErr:  true,
			wantCode: http.StatusForbidden,
		},
		{
			name: "working",
			args: args{
				auth:          authMock.NewMockAuth(nil),
				Authorization: "Bearer mock_token",
				scope:         "scope_1",
				httpMethod:    http.MethodPost,
				reqURLSuffix:  "channels",
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := emptyHandler.Validate(tt.args.auth)
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()
			client := ts.Client()

			reqURL := fmt.Sprintf("%s/%s", ts.URL, tt.args.reqURLSuffix)
			reqURL = strings.TrimSuffix(reqURL, "/")

			req, err := http.NewRequest(tt.args.httpMethod, reqURL, nil)
			if err != nil {
				t.Error("error creating request")
			}

			// adds auth to request header
			req.Header.Add("Authorization", tt.args.Authorization)
			// adds the scope to request header
			req.Header.Add("Scope", tt.args.scope)

			// does request
			res, err := client.Do(req)

			if err != nil {
				t.Error("couldn't receive response")
			}

			got := res.StatusCode
			if !reflect.DeepEqual(got, tt.wantCode) {
				t.Errorf(
					"Handler.Validate() = %v, want %v",
					got,
					tt.wantCode,
				)
			}
		})
	}
}

func Test_getOperation(t *testing.T) {
	type args struct {
		r *http.Request
	}
	createReq := func(method string) *http.Request {
		req, _ := http.NewRequest(
			method,
			"mock-url",
			bytes.NewBuffer([]byte{}),
		)
		return req
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GetOperation",
			args: args{
				r: createReq(http.MethodGet),
			},
			want: "get",
		},
		{
			name: "CreateOperation",
			args: args{
				r: createReq(http.MethodPost),
			},
			want: "create",
		},
		{
			name: "UpdateOperation",
			args: args{
				r: createReq(http.MethodPut),
			},
			want: "update",
		},
		{
			name: "DeleteOperation",
			args: args{
				r: createReq(http.MethodDelete),
			},
			want: "delete",
		},
		{
			name: "InvalidOperation",
			args: args{
				r: createReq(http.MethodPatch),
			},
			want: http.MethodPatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOperation(tt.args.r); got != tt.want {
				t.Errorf("getOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTarget(t *testing.T) {
	type args struct {
		r *http.Request
	}
	const baseURL = "https://localhost:8080/"
	createReq := func(url string) *http.Request {
		req, _ := http.NewRequest(
			http.MethodPost,
			baseURL+url,
			bytes.NewBuffer([]byte{}),
		)
		return req
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple_target",
			args: args{createReq("hello")},
			want: "hello",
		},
		{
			name: "composed_target",
			args: args{createReq("hello/world/")},
			want: "hello/world",
		},
		{
			name: "exception_dapps",
			args: args{createReq("apps")},
			want: "dapp",
		},
		{
			name: "exception_types",
			args: args{createReq("types")},
			want: "type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTarget(tt.args.r); got != tt.want {
				t.Errorf("getTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}
