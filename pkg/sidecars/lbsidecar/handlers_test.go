package lbsidecar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

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

func TestServer_writeMessageHandler(t *testing.T) {

	createMockEnvVars()
	defer deleteMockEnvVars()

	wServer := httptest.NewServer(Init().writeMessageHandler())
	req := http.Client{}

	tests := []struct {
		name    string
		channel string
		msg     models.BrokerMessage
		port    string
		wantErr bool
	}{
		{
			name:    "Channel not listed in 'INSPR_OUTPUT_CHANNELS'",
			channel: "invalidChan1",
			wantErr: true,
		},
		{
			name:    "Env var '<chan>_BROKER' doesn't exist",
			channel: "chan2",
			wantErr: true,
		},
		{
			name:    "Channel avro schema not defined",
			channel: "chan4",
			wantErr: true,
		},
		{
			name:    "Invalid avro schema",
			channel: "chan1",
			wantErr: true,
		},
		{
			name:    "Invalid message given schema",
			channel: "chan3",
			msg: models.BrokerMessage{
				Data: randomStruct{},
			},
			wantErr: true,
		},
		{
			name:    "Invalid request address",
			channel: "chan6",
			msg: models.BrokerMessage{
				Data: "randomMessage",
			},
			wantErr: true,
		},
		{
			name:    "Valid write request",
			channel: "chan5a",
			msg: models.BrokerMessage{
				Data: "randomMessage",
			},
			port: "1107",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testServer *httptest.Server
			if tt.port != "" {
				testServer = createMockedServer(tt.port, tt.channel, nil)
				testServer.Start()
				defer testServer.Close()
			}

			buf, _ := json.Marshal(tt.msg)
			reqInfo, _ := http.NewRequest(http.MethodPost,
				wServer.URL+"/channel/"+tt.channel,
				bytes.NewBuffer(buf))

			resp, err := req.Do(reqInfo)
			if err != nil {
				t.Errorf("Error while doing the request: %v", err)
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

func TestServer_readMessageHandler(t *testing.T) {
	createMockEnvVars()
	defer deleteMockEnvVars()

	rServer := httptest.NewServer(Init().readMessageHandler())
	req := http.Client{}

	tests := []struct {
		name          string
		channel       string
		msg           models.BrokerMessage
		port          string
		setClientPort bool
		wantErr       bool
	}{
		{
			name:    "Channel not listed in 'INSPR_INPUT_CHANNELS'",
			channel: "invalidChan1",
			wantErr: true,
		},

		{
			name:          "Channel avro schema not defined",
			channel:       "chan4",
			wantErr:       true,
			setClientPort: true,
		},
		{
			name:    "Invalid avro schema",
			channel: "chan1",
			wantErr: true,
		},
		{
			name:    "Invalid message given schema",
			channel: "chan3",
			msg: models.BrokerMessage{
				Data: randomStruct{},
			},
			wantErr: true,
		},
		{
			name:    "Invalid scclient request address",
			channel: "chan6",
			msg: models.BrokerMessage{
				Data: "randomMessage",
			},
			wantErr: true,
		},
		{
			name:    "Valid write request",
			channel: "chan5b",
			msg: models.BrokerMessage{
				Data: "randomMessage",
			},
			port:          "1117",
			setClientPort: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setClientPort {
				os.Setenv("INSPR_SCCLIENT_READ_PORT", "1117")
				defer os.Unsetenv("INSPR_SCCLIENT_READ_PORT")
			}
			var testServer *httptest.Server
			if tt.port != "" {
				testServer = createMockedServer(tt.port, tt.channel, tt.msg.Data.(string))
				testServer.Start()
				defer testServer.Close()
			}

			buf, err := encode(tt.channel, tt.msg.Data)
			if err != nil && !tt.wantErr {
				t.Errorf("Unable to encode request: %v", err)
				return
			}
			reqInfo, _ := http.NewRequest(http.MethodPost,
				rServer.URL+"/channel/"+tt.channel,
				bytes.NewBuffer(buf))

			resp, err := req.Do(reqInfo)
			if err != nil {
				t.Errorf("Error while doing the request: %v", err)
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
			name:     "Invalid - port not set",
			msg:      "Hello World!",
			endpoint: "nodename",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.setClientPort {
				os.Setenv("INSPR_SCCLIENT_READ_PORT", "1171")
				defer os.Unsetenv("INSPR_SCCLIENT_READ_PORT")
			}

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

type randomStruct struct{}

func createMockEnvVars() {
	customEnvValues := "chan1@randBroker3;chan3@randBroker4;chan4@randBroker2;chan5a@randBroker1;chan5b@randBroker1;chan6@randBroker5"
	os.Setenv("INSPR_INPUT_CHANNELS", customEnvValues)
	os.Setenv("INSPR_OUTPUT_CHANNELS", customEnvValues)

	os.Setenv("INSPR_LBSIDECAR_WRITE_PORT", "1127")
	os.Setenv("INSPR_LBSIDECAR_READ_PORT", "1137")

	os.Setenv("INSPR_SIDECAR_RANDBROKER2_WRITE_PORT", "somePort1")
	os.Setenv("INSPR_SIDECAR_RANDBROKER2_ADDR", "someAddr1")

	os.Setenv("INSPR_SIDECAR_RANDBROKER3_WRITE_PORT", "somePort1")
	os.Setenv("INSPR_SIDECAR_RANDBROKER3_ADDR", "someAddr1")
	os.Setenv("chan1_SCHEMA", "someSchema")

	os.Setenv("INSPR_SIDECAR_RANDBROKER4_WRITE_PORT", "somePort1")
	os.Setenv("INSPR_SIDECAR_RANDBROKER4_ADDR", "someAddr1")
	os.Setenv("chan3_SCHEMA", `{"type":"string"}`)

	os.Setenv("INSPR_SIDECAR_RANDBROKER5_WRITE_PORT", "somePort1")
	os.Setenv("INSPR_SIDECAR_RANDBROKER5_ADDR", "someAddr1")
	os.Setenv("chan6_SCHEMA", `{"type":"string"}`)

	os.Setenv("INSPR_SIDECAR_RANDBROKER1_WRITE_PORT", "1107")
	os.Setenv("INSPR_SIDECAR_RANDBROKER1_ADDR", "http://localhost")
	os.Setenv("chan5a_SCHEMA", `{"type":"string"}`)
	os.Setenv("chan5a_RESOLVED", "chan5a")
	os.Setenv("chan5b_SCHEMA", `{"type":"string"}`)
	os.Setenv("chan5b_RESOLVED", "chan5b")

}

func deleteMockEnvVars() {
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_LBSIDECAR_WRITE_PORT")
	os.Unsetenv("INSPR_LBSIDECAR_READ_PORT")
	os.Unsetenv("INSPR_SIDECAR_RANDBROKER2_WRITE_PORT")
	os.Unsetenv("INSPR_SIDECAR_RANDBROKER2_ADDR")
	os.Unsetenv("INSPR_SIDECAR_RANDBROKER3_WRITE_PORT")
	os.Unsetenv("INSPR_SIDECAR_RANDBROKER3_ADDR")
	os.Unsetenv("chan1_SCHEMA")
	os.Unsetenv("INSPR_SIDECAR_RANDBROKER4_WRITE_PORT")
	os.Unsetenv("INSPR_SIDECAR_RANDBROKER4_ADDR")
	os.Unsetenv("chan3_SCHEMA")
	os.Unsetenv("INSPR_SIDECAR_RANDBROKER5_WRITE_PORT")
	os.Unsetenv("INSPR_SIDECAR_RANDBROKER5_ADDR")
	os.Unsetenv("chan6_SCHEMA")
	os.Unsetenv("INSPR_SIDECAR_RANDBROKER1_WRITE_PORT")
	os.Unsetenv("INSPR_SIDECAR_RANDBROKER1_ADDR")
	os.Unsetenv("chan5_SCHEMA")
	os.Unsetenv("chan5_RESOLVED")

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
