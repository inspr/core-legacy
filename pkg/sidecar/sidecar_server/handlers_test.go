package sidecarserv

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

func Test_newCustomHandlers(t *testing.T) {
	type args struct {
		server *Server
	}
	tests := []struct {
		name string
		args args
		want *customHandlers
	}{
		{
			name: "successfully_created_custom_handlers",
			args: args{mockServer(nil)},
			want: &customHandlers{Server: mockServer(nil)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCustomHandlers(tt.args.server); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCustomHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}

// sendInRequest is a struct used for all the testing files in this package
// it's contents is a simple { body []byte }
type sendInRequest struct{ body []byte }

// expectedResponse is a struct used for all the testing files in this package
// it's contents is a simple { status int }
type wantedResponse struct{ status int }

func Test_customHandlers_writeMessageHandler(t *testing.T) {
	parsedBody, _ := json.Marshal(models.RequestBody{
		Message: models.Message{},
		Channel: "chan",
	})
	noChanBody, _ := json.Marshal(models.RequestBody{
		Message: models.Message{},
		Channel: "donExist",
	})
	badBody := []byte{0}

	customEnvValues := "chan;testing;banana"
	os.Setenv("INSPR_INPUT_CHANNELS", customEnvValues)
	os.Setenv("INSPR_OUTPUT_CHANNELS", customEnvValues)
	os.Setenv("UNIX_SOCKET_ADDRESS", customEnvValues)

	type args struct {
		send sendInRequest
		want wantedResponse
	}
	tests := []struct {
		name string
		ch   *customHandlers
		args args
	}{
		{
			name: "successful_request",
			ch:   newCustomHandlers(mockServer(nil)),
			args: args{
				send: sendInRequest{parsedBody},
				want: wantedResponse{http.StatusOK},
			},
		},
		{
			name: "unsuccessful_request",
			ch:   newCustomHandlers(mockServer(errors.New("error"))),
			args: args{
				send: sendInRequest{parsedBody},
				want: wantedResponse{http.StatusInternalServerError},
			},
		},
		{
			name: "bad_request",
			ch:   newCustomHandlers(mockServer(nil)),
			args: args{
				send: sendInRequest{badBody},
				want: wantedResponse{http.StatusBadRequest},
			},
		},
		{
			name: "no_channel_request",
			ch:   newCustomHandlers(mockServer(nil)),
			args: args{
				send: sendInRequest{noChanBody},
				want: wantedResponse{http.StatusBadRequest},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// tt.ch.writeMessageHandler(tt.args.w, tt.args.r)
			handlerFunc := http.HandlerFunc(tt.ch.writeMessageHandler)
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			client := ts.Client()
			res, err := client.Post(ts.URL, "", bytes.NewBuffer(tt.args.send.body))
			if err != nil {
				t.Errorf("error making a POST in the httptest server")
			}
			defer res.Body.Close()

			if res.StatusCode != tt.args.want.status {
				t.Errorf("writeMessageHandler = %v, want %v", res, tt.args.want.status)
			}

		})
	}
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("UNIX_SOCKET_ADDRESS")
}

// func Test_customHandlers_readMessageHandler(t *testing.T) {
// 	type args struct {
// 		send sendInRequest
// 		want wantedResponse
// 	}
// 	tests := []struct {
// 		name string
// 		ch   *customHandlers
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.ch.readMessageHandler(tt.args.w, tt.args.r)
// 		})
// 	}
// }

func Test_customHandlers_commitMessageHandler(t *testing.T) {
	parsedBody, _ := json.Marshal(models.RequestBody{
		Message: models.Message{},
		Channel: "chan",
	})
	noChanBody, _ := json.Marshal(models.RequestBody{
		Message: models.Message{},
		Channel: "donExist",
	})
	badBody := []byte{0}

	customEnvValues := "chan;testing;banana"
	os.Setenv("INSPR_INPUT_CHANNELS", customEnvValues)
	os.Setenv("INSPR_OUTPUT_CHANNELS", customEnvValues)
	os.Setenv("UNIX_SOCKET_ADDRESS", customEnvValues)

	type args struct {
		send sendInRequest
		want wantedResponse
	}
	tests := []struct {
		name string
		ch   *customHandlers
		args args
	}{
		{
			name: "successful_request",
			ch:   newCustomHandlers(mockServer(nil)),
			args: args{
				send: sendInRequest{parsedBody},
				want: wantedResponse{http.StatusOK},
			},
		},
		{
			name: "unsuccessful_request",
			ch:   newCustomHandlers(mockServer(errors.New("error"))),
			args: args{
				send: sendInRequest{parsedBody},
				want: wantedResponse{http.StatusInternalServerError},
			},
		},
		{
			name: "bad_request",
			ch:   newCustomHandlers(mockServer(nil)),
			args: args{
				send: sendInRequest{badBody},
				want: wantedResponse{http.StatusBadRequest},
			},
		},
		{
			name: "no_channel_request",
			ch:   newCustomHandlers(mockServer(nil)),
			args: args{
				send: sendInRequest{noChanBody},
				want: wantedResponse{http.StatusBadRequest},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := http.HandlerFunc(tt.ch.commitMessageHandler)
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			client := ts.Client()
			res, err := client.Post(ts.URL, "", bytes.NewBuffer(tt.args.send.body))
			if err != nil {
				t.Errorf("error making a POST in the httptest server")
			}
			defer res.Body.Close()

			if res.StatusCode != tt.args.want.status {
				t.Errorf("writeMessageHandler = %v, want %v", res, tt.args.want.status)
			}

		})
	}
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("UNIX_SOCKET_ADDRESS")
}
