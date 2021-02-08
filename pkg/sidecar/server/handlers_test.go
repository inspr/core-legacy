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

	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
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
// the models.RequestBody that sets a standard struct to be sent in each
// request made in the handler func.
func generateTestCases() []testCaseStruct {
	// default values used in the test cases
	parsedBody, _ := json.Marshal(models.RequestBody{
		Message: models.BodyMessage{Data: "data"},
		Channel: "chan",
	})
	noChanBody, _ := json.Marshal(models.RequestBody{
		Message: models.BodyMessage{Data: "data"},
		Channel: "donExist",
	})
	badBody := []byte{0}

	return []testCaseStruct{
		{
			name: "successful_request",
			ch:   newCustomHandlers(MockServer(nil)),
			args: args{
				send: sendInRequest{parsedBody},
				want: wantedResponse{http.StatusOK},
			},
		},
		{
			name: "unsuccessful_request",
			ch:   newCustomHandlers(MockServer(errors.New("error"))),
			args: args{
				send: sendInRequest{parsedBody},
				want: wantedResponse{http.StatusInternalServerError},
			},
		},
		{
			name: "bad_request",
			ch:   newCustomHandlers(MockServer(nil)),
			args: args{
				send: sendInRequest{badBody},
				want: wantedResponse{http.StatusBadRequest},
			},
		},
		{
			name: "no_channel_request",
			ch:   newCustomHandlers(MockServer(nil)),
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
	os.Setenv("UNIX_SOCKET_ADDRESS", unixSocketAddr)
}

// deleteMockEnvVars - deletes the env values used in the tests functions
func deleteMockEnvVars() {
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("UNIX_SOCKET_ADDRESS")
}

func Test_newCustomHandlers(t *testing.T) {
	createMockEnvVars()
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
				Server:    MockServer(nil),
				insprVars: environment.GetEnvironment(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCustomHandlers(tt.args.server); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCustomHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
	deleteMockEnvVars()
}

func Test_customHandlers_writeMessageHandler(t *testing.T) {
	createMockEnvVars()
	tests := generateTestCases()

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
	deleteMockEnvVars()
}

func Test_customHandlers_readMessageHandler(t *testing.T) {
	createMockEnvVars()
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
				msg := models.BrokerResponse{}
				err := json.NewDecoder(res.Body).Decode(&msg)

				// if it failed to parse body
				if err != nil {
					t.Log("Failed to parse the receive message")
					t.Errorf("readMessageHandler.Body error = %v, want %v", err, nil)
				}

				// if channel isn't 'chan'
				expectedData, _ := MockServer(nil).Reader.ReadMessage("")
				if msg.Data != expectedData.Data {
					t.Errorf("readMessageHandler.Body error, field 'data' = %v, want %v", msg.Data, expectedData)
				}
			}
		})
	}
	deleteMockEnvVars()
}

func Test_customHandlers_commitMessageHandler(t *testing.T) {
	createMockEnvVars()

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
	deleteMockEnvVars()
}

func Test_handlers_existsInSlice(t *testing.T) {
	tests := []struct {
		testName    string
		channel     string
		channelList []string
		want        bool
	}{
		{
			testName:    "found_channel",
			channel:     "a",
			channelList: []string{"a", "b", "c", "d"},
			want:        true,
		},
		{
			testName:    "no_channel_found",
			channel:     "non-existant",
			channelList: []string{"a", "b", "c", "d"},
			want:        false,
		},
		{
			testName:    "empty_channelList",
			channel:     "empty",
			channelList: []string{},
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got := existsInSlice(tt.channel, tt.channelList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCustomHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}
