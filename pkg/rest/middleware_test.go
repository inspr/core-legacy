package rest

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	authentication "gitlab.inspr.dev/inspr/core/cmd/insprd/auth"
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
		auth authentication.Auth
	}
	tests := []struct {
		name string
		h    Handler
		args args
		want Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Validate(tt.args.auth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
