package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/ierrors"
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
			err:  errors.New(""),
			want: http.StatusInternalServerError,
		},
		{
			name: "InsprErrors_NotFound",
			err:  ierrors.New("").NotFound(),
			want: http.StatusNotFound,
		},
		{
			name: "InsprErrors_AlreadyExists",
			err:  ierrors.New("").AlreadyExists(),
			want: http.StatusConflict,
		},
		{
			name: "InsprErrors_InternalServer",
			err:  ierrors.New("").InternalServer(),
			want: http.StatusInternalServerError,
		},
		{
			name: "InsprErrors_InvalidName",
			err:  ierrors.New("").InvalidName(),
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_InvalidApp",
			err:  ierrors.New("").InvalidApp(),
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_InvalidChannel",
			err:  ierrors.New("").InvalidChannel(),
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_InvalidType",
			err:  ierrors.New("").InvalidType(),
			want: http.StatusForbidden,
		},
		{
			name: "InsprErrors_BadRequest",
			err:  ierrors.New("").BadRequest(),
			want: http.StatusBadRequest,
		},
		{
			name: "InsprErrors_Unknown_ErrCode",
			err:  ierrors.New(""),
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

			// TODO REVIEW
			errorMessage := ierrors.New("")
			json.Unmarshal(rr.Body.Bytes(), &errorMessage)

			gotCode := ierrors.Code(errorMessage)
			wantCode := ierrors.Code(tt.err)

			if !reflect.DeepEqual(gotCode, wantCode) {
				t.Errorf("JSON(w,code,data)=%v, want %v", gotCode, wantCode)
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

	got := ierrors.New("")
	json.NewDecoder(resp.Body).Decode(&got)

	want := ierrors.New("This is a panic error").InternalServer()

	if !reflect.DeepEqual(want.Error(), got.Error()) {
		t.Errorf("RecoverFromPanic=%v, want %v", got, want)
	}

}

func TestUnmarshalERROR(t *testing.T) {
	type args struct {
		r io.Reader
	}

	generateBody := func(body string) io.Reader {
		bodyBytes, _ := json.Marshal(body)
		return bytes.NewBuffer(bodyBytes)
	}
	generateErrBody := func(err error) io.Reader {
		errBytes, _ := json.Marshal(err)
		return bytes.NewBuffer(errBytes)
	}

	tests := []struct {
		name     string
		args     args
		want     error
		wantCode ierrors.ErrCode
	}{
		{
			name: "basic_unmarshal_error",
			args: args{r: generateErrBody(
				ierrors.New("no permission to create dapp").Forbidden(),
			)},
			want:     ierrors.New("no permission to create dapp").Forbidden(),
			wantCode: ierrors.Forbidden,
		},
		{
			name:     "basic_unmarshal_empty_error",
			args:     args{r: generateBody("nothing")},
			want:     defaultErr,
			wantCode: ierrors.Unknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.r)
			err := UnmarshalERROR(tt.args.r)
			if err.Error() != tt.want.Error() {
				t.Errorf("UnmarshalERROR() error = %v, wantErr %v",
					err,
					tt.want,
				)
			}
		})
	}
}
