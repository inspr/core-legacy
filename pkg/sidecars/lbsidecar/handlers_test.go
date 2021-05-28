package lbsidecar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/sidecars/models"
)

func createMockedServer(port, ch string, msg interface{}) *httptest.Server {
	// create a listener with the desired port.
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewUnstartedServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			channel := strings.TrimPrefix(r.URL.Path, "/")

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
			channel: "chan5",
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
				wServer.URL+"/"+tt.channel,
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
				fmt.Printf("Response body: %v", resp.Body)
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
			name:          "Invalid avro schema",
			channel:       "chan1",
			wantErr:       true,
			setClientPort: true,
		},
		{
			name:    "Invalid message given schema",
			channel: "chan3",
			msg: models.BrokerMessage{
				Data: randomStruct{},
			},
			wantErr:       true,
			setClientPort: true,
		},
		{
			name:    "Invalid request address",
			channel: "chan6",
			msg: models.BrokerMessage{
				Data: "randomMessage",
			},
			wantErr:       true,
			setClientPort: true,
		},
		{
			name:    "Valid write request",
			channel: "chan5",
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
				rServer.URL+"/"+tt.channel,
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

type randomStruct struct{}

func createMockEnvVars() {
	customEnvValues := "chan1_randBroker3;chan3_randBroker4;chan4_randBroker2;chan5_randBroker1;chan6_randBroker5"
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
	os.Setenv("chan5_SCHEMA", `{"type":"string"}`)

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
}
