package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	authentication "gitlab.inspr.dev/inspr/core/cmd/insprd/auth"
	authMock "gitlab.inspr.dev/inspr/core/cmd/insprd/auth/mocks"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
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
	type args struct {
		auth        authentication.Auth
		headerValue string
		scope       string
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
				auth: authMock.NewMockAuth(nil),
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "invalid_token",
			args: args{
				auth: authMock.NewMockAuth(
					&ierrors.InsprError{Code: ierrors.InvalidToken},
				),
				headerValue: "Bearer mock_token",
			},
			wantErr:  true,
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "expired_token",
			args: args{
				auth: authMock.NewMockAuth(
					&ierrors.InsprError{Code: ierrors.ExpiredToken},
				),
				headerValue: "Bearer mock_token",
			},
			wantErr:  true,
			wantCode: http.StatusOK,
		},
		{
			name: "unknown_error",
			args: args{
				auth: authMock.NewMockAuth(
					errors.New("mock_error"),
				),
				headerValue: "Bearer mock_token",
			},
			wantErr:  true,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "scope_error",
			args: args{
				auth:        authMock.NewMockAuth(nil),
				headerValue: "Bearer mock_token",
				scope:       "no_valid_scope",
			},
			wantErr:  true,
			wantCode: http.StatusForbidden,
		},
		{
			name: "scope_error_doesnt_have_prefix",
			args: args{
				auth:        authMock.NewMockAuth(nil),
				headerValue: "Bearer mock_token",
				scope:       "wrongScope.scope_1",
			},
			wantErr:  true,
			wantCode: http.StatusForbidden,
		},
		{
			name: "working",
			args: args{
				auth:        authMock.NewMockAuth(nil),
				headerValue: "Bearer mock_token",
				scope:       "scope_1",
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

			req, err := http.NewRequest(http.MethodPost, ts.URL, nil)
			if err != nil {
				t.Error("error creating request")
			}

			// adds auth to request header
			req.Header.Add("Authorization", tt.args.headerValue)

			// scope to body
			scopeData := struct {
				Scope string `json:"scope"`
			}{
				Scope: tt.args.scope,
			}
			scopeBytes, _ := json.Marshal(scopeData)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(scopeBytes))

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
