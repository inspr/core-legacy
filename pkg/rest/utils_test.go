package rest

import (
	"encoding/json"
	"errors"
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
	rr := httptest.NewRecorder()
	type args struct {
		w   http.ResponseWriter
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "non_InsprErrors",
			args: args{
				w:   rr,
				err: errors.New("server crashed"),
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "InsprErrors_NotFound",
			args: args{
				w:   rr,
				err: ierrors.NewError().NotFound().Build(),
			},
			want: http.StatusNotFound,
		},
		{
			name: "InsprErrors_AlreadyExists",
			args: args{
				w:   rr,
				err: ierrors.NewError().AlreadyExists().Build(),
			},
			want: http.StatusConflict,
		},
		{
			name: "InsprErrors_InternalServer",
			args: args{
				w:   rr,
				err: ierrors.NewError().InternalServer().Build(),
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "InsprErrors_InvalidName",
			args: args{
				w:   rr,
				err: ierrors.NewError().InvalidName().Build(),
			},
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_InvalidApp",
			args: args{
				w:   rr,
				err: ierrors.NewError().InvalidApp().Build(),
			},
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_InvalidChannel",
			args: args{
				w:   rr,
				err: ierrors.NewError().InvalidChannel().Build(),
			},
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_InvalidChannelType",
			args: args{
				w:   rr,
				err: ierrors.NewError().InvalidChannelType().Build(),
			},
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_BadRequest",
			args: args{
				w:   rr,
				err: ierrors.NewError().BadRequest().Build(),
			},
			want: http.StatusBadRequest,
		},
		{
			name: "InsprErrors_Unknown_ErrCode",
			args: args{
				w:   rr,
				err: &ierrors.InsprError{Code: 9999},
			},
			want: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ERROR(tt.args.w, tt.args.err)
			if status := rr.Result().StatusCode; status != tt.want {
				t.Errorf("JSON(w,code,data)=%v, want %v", status, tt.want)
			}

			errorMessage := ierrors.InsprError{}

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
