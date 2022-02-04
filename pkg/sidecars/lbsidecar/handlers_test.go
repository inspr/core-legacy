package lbsidecar

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"inspr.dev/inspr/pkg/rest"
	"inspr.dev/inspr/pkg/sidecars/models"
)

func createMockedServer(port, ch string, msg interface{}) *httptest.Server {
	// create a listener with the desired port.
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewUnstartedServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			channel := strings.TrimPrefix(r.URL.Path, "/channel/")

			var receivedData models.BrokerMessage
			json.NewDecoder(r.Body).Decode(&receivedData)

			if (ch != channel) || (msg != receivedData.Data) {
				rest.ERROR(w, fmt.Errorf("invalid channel or message"))
				return
			}

			rest.JSON(w, http.StatusOK, nil)
		}),
	)
	// NewUnstartedServer creates a listener. Close that listener and replace
	// with the one we created.
	ts.Listener.Close()
	ts.Listener = l

	return ts
}

func createRouteMockedServer(port string, expectedMsg interface{}) *httptest.Server {
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}

	mockServer := httptest.NewUnstartedServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				rest.ERROR(w, fmt.Errorf("error while reading msg body"))
				return
			}

			fmt.Printf("receivedData = %v\n", string(body))

			if expectedMsg != string(body) {
				rest.ERROR(w, fmt.Errorf("invalid message"))
				return
			}

			rest.JSON(w, http.StatusOK, nil)
		}),
	)

	mockServer.Listener.Close()
	mockServer.Listener = listener

	return mockServer

}

type mockWriter struct {
	writeMessage func(channel string, message []byte) error
}

func (m *mockWriter) WriteMessage(channel string, message []byte) error {
	return m.writeMessage(channel, message)
}

func (m *mockWriter) Close() {}

func (m *mockWriter) Producer() *kafka.Producer { return nil }

func TestServer_writeMessageHandler(t *testing.T) {
	createMockEnvVars()
	defer deleteMockEnvVars()
	tests := []struct {
		readerFunc func(t *testing.T) models.Reader
		writerFunc func(t *testing.T) models.Writer
		channel    string
		name       string
		wantErr    bool
		message    models.BrokerMessage
	}{
		{
			name:    "correct behaviour test",
			channel: "chan",
			message: models.BrokerMessage{
				Data: "randomMessage",
			},
		},
		{
			name:    "invalid channel",
			channel: "invalid",
			message: models.BrokerMessage{
				Data: "randomMessage",
			},
			wantErr: true,
		},
		{
			name:    "invalid data for marshalling",
			channel: "chan2",
			message: models.BrokerMessage{
				Data: 45,
			},
			wantErr: true,
		},
		{
			name:    "invalid broker response",
			channel: "chan3",
			message: models.BrokerMessage{
				Data: "randomMessage",
			},
			wantErr: true,
			writerFunc: func(t *testing.T) models.Writer {
				return &mockWriter{
					writeMessage: func(channel string, message []byte) error {
						return errors.New("this is an error")
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.writerFunc == nil {
				tt.writerFunc = func(t *testing.T) models.Writer {
					return &mockWriter{
						writeMessage: func(channel string, message []byte) error {
							// this should check two things, if the test channel was resolved correctly
							// and if the test message was correctly encoded

							if channel != tt.channel {
								t.Errorf(
									"Server_writeMessageHandler channel resolved to = %s, expected = %s",
									channel,
									tt.channel,
								)
							}

							resolvedCh, _ := getResolvedChannel(channel)

							bMessage, err := readMessage(resolvedCh, message)
							if err != nil {
								t.Error(err)
							}

							if !reflect.DeepEqual(bMessage, tt.message) {
								t.Errorf(
									"Server_writeMessageHandler message = %v, want = %v",
									bMessage,
									tt.message,
								)
							}

							return nil
						},
					}
				}
			}
			bh := models.NewBrokerHandler("someBroker", nil, tt.writerFunc(t))
			s := Init(bh)
			server := httptest.NewServer(s.writeMessageHandler())
			defer server.Close()
			client := &http.Client{}

			buf, _ := json.Marshal(tt.message)
			resp, err := client.Post(
				fmt.Sprintf("%s/channel/%s", server.URL, tt.channel),
				"application/octet-stream",
				bytes.NewBuffer(buf),
			)
			if err != nil || resp.StatusCode != http.StatusOK {
				err = rest.UnmarshalERROR(resp.Body)
				if (err != nil) != tt.wantErr {
					t.Errorf("Server_writeMessageHandler err = %v, wantErr = %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestServer_routeReceiveHandler(t *testing.T) {

	createMockEnvVars()
	defer deleteMockEnvVars()

	rServer := httptest.NewServer(Init().routeReceiveHandler())
	req := http.Client{}

	tests := []struct {
		name          string
		msg           string
		endpoint      string
		port          string
		setClientPort bool
		wantErr       bool
		want          rest.Handler
	}{
		{
			name:          "Valid route request",
			msg:           "Hello World!",
			endpoint:      "nodename/hello",
			port:          "1171",
			setClientPort: true,
		},
		{
			name:          "Invalid - port not set",
			msg:           "Hello World!",
			endpoint:      "nodename",
			wantErr:       true,
			setClientPort: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var testServer *httptest.Server
			if tt.port != "" {
				testServer = createRouteMockedServer(tt.port, tt.msg)
				testServer.Start()
				defer testServer.Close()
			}

			reqInfo, _ := http.NewRequest(http.MethodPost, rServer.URL+"/route/"+tt.endpoint, bytes.NewBuffer([]byte(tt.msg)))

			resp, err := req.Do(reqInfo)
			if err != nil {
				t.Errorf("Error while doing the request: %v", err)
				return
			}

			if resp.StatusCode != http.StatusOK && !tt.wantErr {
				t.Errorf("Received status %v, wanted %v", resp.StatusCode, http.StatusOK)
				return
			}

		})
	}
}

func createMockEnvVars() {
	customEnvValues := "chan@someBroker;testing@someBroker;chan2@someBroker"
	os.Setenv("INSPR_INPUT_CHANNELS", customEnvValues)
	os.Setenv("INSPR_OUTPUT_CHANNELS", customEnvValues)

	os.Setenv("INSPR_LBSIDECAR_WRITE_PORT", "1127")
	os.Setenv("INSPR_LBSIDECAR_READ_PORT", "1137")
	os.Setenv("INSPR_SCCLIENT_READ_PORT", "1171")

	os.Setenv("chan_RESOLVED", "someTopic")
	os.Setenv("testing_RESOLVED", "someTopic")
	os.Setenv("chan2_RESOLVED", "someTopic")

	os.Setenv("someTopic_SCHEMA", `{"type":"string"}`)

}

func deleteMockEnvVars() {
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")

	os.Unsetenv("INSPR_LBSIDECAR_WRITE_PORT")
	os.Unsetenv("INSPR_LBSIDECAR_READ_PORT")
	os.Unsetenv("INSPR_SCCLIENT_READ_PORT")

	os.Unsetenv("chan_RESOLVED")
	os.Unsetenv("testing_RESOLVED")
	os.Unsetenv("chan2_RESOLVED")

	os.Unsetenv("someTopic_SCHEMA")
}

func TestServer_sendRequest(t *testing.T) {
	createMockEnvVars()
	defer deleteMockEnvVars()

	wServer := httptest.NewServer(Init().sendRouteRequest())
	client := http.Client{}

	tests := []struct {
		name    string
		path    string
		msg     string
		route   string
		data    string
		wantErr bool
	}{
		{
			name:    "valid route request",
			path:    "rt1/edp1",
			msg:     "mock_string",
			route:   "rt1_ROUTE",
			data:    "http://localhost:3301;edp1;edp2",
			wantErr: false,
		},
		{
			name:    "invalid route request",
			path:    "rt1/edp1",
			msg:     "mock_string",
			route:   "non_ROUTE",
			data:    "http://localhost:3301;edp1;edp2",
			wantErr: true,
		},
		{
			name:    "invalid route endpoint",
			path:    "rt1/edp3",
			msg:     "mock_string",
			route:   "rt1_ROUTE",
			data:    "http://localhost:3301;edp1;edp2",
			wantErr: true,
		},
		{
			name:    "invalid route address",
			path:    "rt1/edp1",
			msg:     "mock_string",
			route:   "rt1_ROUTE",
			data:    "localhost:3301;edp1;edp2",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.route, tt.data)
			defer os.Unsetenv(tt.route)

			testServer := createMockedRoute("3301", tt.msg)
			testServer.Start()
			defer testServer.Close()

			buf, _ := json.Marshal(tt.msg)
			reqInfo, _ := http.NewRequest(http.MethodPost,
				wServer.URL+"/route/"+tt.path,
				bytes.NewBuffer(buf))

			resp, err := client.Do(reqInfo)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Error while doing the request: %v", err)
				}
				return
			}

			if tt.wantErr && (resp.StatusCode == http.StatusOK) {
				t.Errorf("Wanted error, received 'nil'")
				return
			}

			if !tt.wantErr && (resp.StatusCode != http.StatusOK) {
				t.Errorf("Received status %v, wanted %v", resp.StatusCode, http.StatusOK)
				return
			}
		})
	}
}

func createMockedRoute(port string, expectedMsg interface{}) *httptest.Server {
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}

	mockServer := httptest.NewUnstartedServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			var receivedData string
			err := json.NewDecoder(r.Body).Decode(&receivedData)
			if err != nil {
				rest.ERROR(w, fmt.Errorf("error while reading msg body"))
				return
			}

			if expectedMsg != receivedData {
				rest.ERROR(w, fmt.Errorf("invalid message"))
				return
			}

			rest.JSON(w, http.StatusOK, nil)
		}),
	)

	mockServer.Listener.Close()
	mockServer.Listener = listener

	return mockServer

}
