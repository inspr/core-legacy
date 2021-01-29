package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	type args struct {
		w          http.ResponseWriter
		statusCode int
		data       interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success_test",
			args: args{
				w:          rr,
				statusCode: 200,
				data:       "testing",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			JSON(tt.args.w, tt.args.statusCode, tt.args.data)
			if status := rr.Result().StatusCode; status != tt.args.statusCode {
				t.Errorf("JSON(w,code,data)=%v, want %v", status, tt.args.statusCode)
			}
			decodedData, _ := json.Marshal(tt.args.data)
			bodyData, _ := ioutil.ReadAll(rr.Body)
			bodyData = bodyData[:len(bodyData)-1] // removing EOF byte

			if !reflect.DeepEqual(bodyData, decodedData) {
				t.Errorf("JSON(w,code,data)=%v, want %v", bodyData, decodedData)
			}
		})
	}
}

func TestERROR(t *testing.T) {
	rr := httptest.NewRecorder()
	type args struct {
		w          http.ResponseWriter
		statusCode int
		err        error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful_test",
			args: args{
				w:          rr,
				statusCode: http.StatusBadRequest,
				err:        errors.New("My testing error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ERROR(tt.args.w, tt.args.statusCode, tt.args.err)
			if status := rr.Result().StatusCode; status != tt.args.statusCode {
				t.Errorf("JSON(w,code,data)=%v, want %v", status, tt.args.statusCode)
			}

			var errorMessage struct {
				Error string `json:"error"`
			}

			err := json.Unmarshal(rr.Body.Bytes(), &errorMessage)
			if err != nil {
				t.Fatal("Failed to parse the body bytes from the request")
			}

			if !reflect.DeepEqual(errorMessage.Error, tt.args.err.Error()) {
				t.Errorf("JSON(w,code,data)=%v, want %v", errorMessage.Error, tt.args.err.Error())
			}
		})
	}
}
