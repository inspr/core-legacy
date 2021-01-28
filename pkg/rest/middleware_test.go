package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
