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

	env "inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/rest"
	"inspr.dev/inspr/pkg/sidecar/models"
)

// sendInRequest is a struct used for all the testing files in this package
// it's contents is a simple { body []byte }
type sendInRequest struct{ body []byte }

// expectedResponse is a struct used for all the testing files in this package
// it's contents is a simple { status int }
type wantedResponse struct{ status int }

// args is a struct to setup the values to be sent and the ones
// that are expected in the request
type args struct {
	send sendInRequest
	want wantedResponse
}

// testCaseStruct stores the complete struct used in testing
// containing all the structs defined in the testing file
// sets the name of the test, it's handler and it's args
type testCaseStruct struct {
	name string
	ch   *customHandlers
	args args
}

// generateTestCases returns the tests cases values to be used in each
// handle test, the reason for them to share tests cases is because of
// the models.BrokerData that sets a standard struct to be sent in each
// request made in the handler func.
func generateTestCases() []testCaseStruct {
	// default values used in the test cases
	parsedBody, _ := json.Marshal(models.BrokerData{
		Message: models.Message{Data: "data"},
		Channel: "chan",
	})
	noChanBody, _ := json.Marshal(models.BrokerData{
		Message: models.Message{Data: "data"},
		Channel: "donExist",
	})
	badBody := []byte{0}

	// constants used in the tests
	normalCustomHandler := newCustomHandlers(&MockServer(nil).Mutex, MockServer(nil).Reader, MockServer(nil).Writer)
	err := errors.New("error")
	throwCustomHandler := newCustomHandlers(&MockServer(err).Mutex, MockServer(err).Reader, MockServer(err).Writer)

	return []testCaseStruct{
		{
			name: "successful_request",
			ch:   normalCustomHandler,
			args: args{
				send: sendInRequest{parsedBody},
				want: wantedResponse{http.StatusOK},
			},
		},
		{
			name: "unsuccessful_request",
			ch:   throwCustomHandler,
			args: args{
				send: sendInRequest{parsedBody},
				want: wantedResponse{http.StatusInternalServerError},
			},
		},
		{
			name: "bad_request",
			ch:   normalCustomHandler,
			args: args{
				send: sendInRequest{badBody},
				want: wantedResponse{http.StatusBadRequest},
			},
		},
		{
			name: "no_channel_request",
			ch:   normalCustomHandler,
			args: args{
				send: sendInRequest{noChanBody},
				want: wantedResponse{http.StatusBadRequest},
			},
		},
	}
}

// createMockEnvVars - sets up the env values to be used in the tests functions
func createMockEnvVars() {
	customEnvValues := "chan;testing;banana"
	var unixSocketAddr = "/tmp/insprd.sock"
	os.Setenv("INSPR_INPUT_CHANNELS", customEnvValues)
	os.Setenv("INSPR_OUTPUT_CHANNELS", customEnvValues)
	os.Setenv("INSPR_UNIX_SOCKET", unixSocketAddr)
	os.Setenv("INSPR_APP_CTX", "random.ctx")
	os.Setenv("INSPR_ENV", "test")
	os.Setenv("INSPR_APP_ID", "appid")
}

// deleteMockEnvVars - deletes the env values used in the tests functions
func deleteMockEnvVars() {
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("INSPR_UNIX_SOCKET")
	os.Unsetenv("INSPR_APP_CTX")
	os.Unsetenv("INSPR_ENV")
	os.Unsetenv("INSPR_APP_ID")
}

func Test_newCustomHandlers(t *testing.T) {
	env.SetMockEnv()
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
			args: args{MockServer(nil)},
			want: &customHandlers{
				Locker:         &MockServer(nil).Mutex,
				r:              MockServer(nil).Reader,
				w:              MockServer(nil).Writer,
				InputChannels:  env.GetInputChannels(),
				OutputChannels: env.GetOutputChannels(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newCustomHandlers(&MockServer(nil).Mutex, MockServer(nil).Reader, MockServer(nil).Writer)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCustomHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
	env.UnsetMockEnv()
}

func Test_customHandlers_writeMessageHandler(t *testing.T) {
	env.SetMockEnv()
	tests := generateTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	env.UnsetMockEnv()
}

func Test_customHandlers_readMessageHandler(t *testing.T) {
	env.SetMockEnv()
	tests := generateTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := http.HandlerFunc(tt.ch.readMessageHandler)
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			client := ts.Client()
			res, err := client.Post(ts.URL, "", bytes.NewBuffer(tt.args.send.body))
			if err != nil {
				t.Errorf("error making a POST in the httptest server")
			}
			defer res.Body.Close()

			// identifies if it responded properly
			if res.StatusCode != tt.args.want.status {
				t.Errorf("readMessageHandler = %v, want %v", res, tt.args.want.status)
			}

			if res.StatusCode != http.StatusOK { // reading error
				err := rest.UnmarshalERROR(res.Body)
				if err == nil {
					t.Errorf("readMessageHandler.Body error = %v, want 'nil'", err)
				}
			} else { //reading message

				// reads response and checks for the default mock values
				msg := models.BrokerData{}
				err := json.NewDecoder(res.Body).Decode(&msg)

				// if it failed to parse body
				if err != nil {
					t.Log("Failed to parse the receive message")
					t.Errorf("readMessageHandler.Body error = %v, want %v", err, nil)
				}

				// if channel isn't 'chan'
				expectedData, _ := MockServer(nil).Reader.ReadMessage(msg.Channel)
				if msg.Message.Data != expectedData.Message.Data {
					t.Errorf("readMessageHandler.Body error, field 'data' = %v, want %v", msg.Message.Data, expectedData)
				}
			}
		})
	}
	env.UnsetMockEnv()
}

func Test_customHandlers_commitMessageHandler(t *testing.T) {
	env.SetMockEnv()

	tests := generateTestCases()
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
	env.UnsetMockEnv()
}
