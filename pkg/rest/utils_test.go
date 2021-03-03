package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
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
				statusCode: http.StatusOK,
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
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "non_InsprErrors",
			err:  errors.New("server crashed"),
			want: http.StatusInternalServerError,
		},
		{
			name: "InsprErrors_NotFound",
			err:  ierrors.NewError().NotFound().Build(),
			want: http.StatusNotFound,
		},
		{
			name: "InsprErrors_AlreadyExists",
			err:  ierrors.NewError().AlreadyExists().Build(),
			want: http.StatusConflict,
		},
		{
			name: "InsprErrors_InternalServer",
			err:  ierrors.NewError().InternalServer().Build(),
			want: http.StatusInternalServerError,
		},
		{
			name: "InsprErrors_InvalidName",
			err:  ierrors.NewError().InvalidName().Build(),
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_InvalidApp",
			err:  ierrors.NewError().InvalidApp().Build(),
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_InvalidChannel",
			err:  ierrors.NewError().InvalidChannel().Build(),
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_InvalidChannelType",
			err:  ierrors.NewError().InvalidChannelType().Build(),
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_BadRequest",
			err:  ierrors.NewError().BadRequest().Build(),
			want: http.StatusBadRequest,
		},
		{
			name: "InsprErrors_Unknown_ErrCode",
			err:  &ierrors.InsprError{Code: 9999},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		rr := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			ERROR(rr, tt.err)
			if status := rr.Result().StatusCode; status != tt.want {
				t.Errorf("JSON(w,code,data)=%v, want %v", status, tt.want)
			}
			var errorMessage ierrors.InsprError
			json.Unmarshal(rr.Body.Bytes(), &errorMessage)

			if !reflect.DeepEqual(errorMessage.Message, tt.err.Error()) {
				t.Errorf("JSON(w,code,data)=%v, want %v", errorMessage.Message, tt.err.Error())
			}
		})
	}
}

func TestRecoverFromPanic(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer RecoverFromPanic(w)
		panic("This is a panic error")
	}))
	defer ts.Close()

	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println(err)
	}

	var got *ierrors.InsprError
	json.NewDecoder(resp.Body).Decode(&got)

	want := ierrors.NewError().InternalServer().Message("This is a panic error").Build()

	if !reflect.DeepEqual(want.Message, got.Message) || !reflect.DeepEqual(want.Code, got.Code) {
		t.Errorf("RecoverFromPanic=%v, want %v", got, want)
	}

}

func callRecoverPanic(w http.ResponseWriter) {
	defer RecoverFromPanic(w)
	panic("Panic from Test")
}
