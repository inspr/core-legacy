/*
Package rest contains the functions
that make it easier to manager api
handler functions
*/
package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/ierrors"
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
			name: "conversion testing",
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

func TestHandler_Recover(t *testing.T) {
	var manipulation = func(h Handler) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			h.Recover(nil)(w, r)
		}
	}

	var panicHandler Handler = func(w http.ResponseWriter, r *http.Request) {
		panic("Panic Test")
	}

	req, err := http.NewRequest(http.MethodGet, "/testing", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	successHandler := http.HandlerFunc(manipulation(panicHandler))
	successHandler.ServeHTTP(rr, req)

	body := rr.Result().Body

	var got *ierrors.InsprError
	json.NewDecoder(body).Decode(&got)

	want := ierrors.NewError().InternalServer().Message("Panic Test").Build()

	if !reflect.DeepEqual(want.Message, got.Message) || !reflect.DeepEqual(want.Code, got.Code) {
		t.Errorf("RecoverFromPanic=%v, want %v", got, want)
	}
}
