/*
Package rest contains the functions
that make it easier to manager api
handler functions
*/
package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler_HTTPHandlerFunc(t *testing.T) {
	var myHandler Handler = func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("printing something")
	}
	var HTTPMyHandler http.HandlerFunc = http.HandlerFunc(myHandler)
	tests := []struct {
		name string
		h    Handler
		want http.HandlerFunc
	}{
		{
			name: "convertion testing",
			h:    myHandler,
			want: HTTPMyHandler,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.HTTPHandlerFunc(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Handler.HTTPHandlerFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Get(t *testing.T) {
	// manipulation - calls .Put() of the handler in the parameter
	var manipulation = func(h Handler) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			h.Get()(w, r)
		}
	}
	tests := []struct {
		name          string
		method        string
		customHandler Handler
		want          int
	}{
		{
			name:          "success_method",
			method:        http.MethodGet,
			customHandler: manipulation(emptyHandler),
			want:          http.StatusOK,
		},
		{
			name:          "fail_method",
			method:        http.MethodPut,
			customHandler: manipulation(emptyHandler),
			want:          http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		// sets up the test server
		req, err := http.NewRequest(tt.method, "/testing", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		successHandler := http.HandlerFunc(tt.customHandler)
		successHandler.ServeHTTP(rr, req)

		if status := rr.Result().StatusCode; status != tt.want {
			t.Errorf("Handler.JSON() = %v, want %v", status, tt.want)
		}
	}
}

func TestHandler_Post(t *testing.T) {
	// manipulation - calls .Put() of the handler in the parameter
	var manipulation = func(h Handler) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			h.Post()(w, r)
		}
	}
	tests := []struct {
		name          string
		method        string
		customHandler Handler
		want          int
	}{
		{
			name:          "success_method",
			method:        http.MethodPost,
			customHandler: manipulation(emptyHandler),
			want:          http.StatusOK,
		},
		{
			name:          "fail_method",
			method:        http.MethodGet,
			customHandler: manipulation(emptyHandler),
			want:          http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		// sets up the test server
		req, err := http.NewRequest(tt.method, "/testing", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		successHandler := http.HandlerFunc(tt.customHandler)
		successHandler.ServeHTTP(rr, req)

		if status := rr.Result().StatusCode; status != tt.want {
			t.Errorf("Handler.JSON() = %v, want %v", status, tt.want)
		}
	}
}

func TestHandler_Delete(t *testing.T) {
	// manipulation - calls .Put() of the handler in the parameter
	var manipulation = func(h Handler) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			h.Delete()(w, r)
		}
	}
	tests := []struct {
		name          string
		method        string
		customHandler Handler
		want          int
	}{
		{
			name:          "success_method",
			method:        http.MethodDelete,
			customHandler: manipulation(emptyHandler),
			want:          http.StatusOK,
		},
		{
			name:          "fail_method",
			method:        http.MethodGet,
			customHandler: manipulation(emptyHandler),
			want:          http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		// sets up the test server
		req, err := http.NewRequest(tt.method, "/testing", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		successHandler := http.HandlerFunc(tt.customHandler)
		successHandler.ServeHTTP(rr, req)

		if status := rr.Result().StatusCode; status != tt.want {
			t.Errorf("Handler.JSON() = %v, want %v", status, tt.want)
		}
	}
}

func TestHandler_Put(t *testing.T) {
	// manipulation - calls .Put() of the handler in the parameter
	var manipulation = func(h Handler) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			h.Put()(w, r)
		}
	}
	tests := []struct {
		name          string
		method        string
		customHandler Handler
		want          int
	}{
		{
			name:          "success_method",
			method:        http.MethodPut,
			customHandler: manipulation(emptyHandler),
			want:          http.StatusOK,
		},
		{
			name:          "fail_method",
			method:        http.MethodGet,
			customHandler: manipulation(emptyHandler),
			want:          http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		// sets up the test server
		req, err := http.NewRequest(tt.method, "/testing", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		successHandler := http.HandlerFunc(tt.customHandler)
		successHandler.ServeHTTP(rr, req)

		if status := rr.Result().StatusCode; status != tt.want {
			t.Errorf("Handler.JSON() = %v, want %v", status, tt.want)
		}
	}
}
