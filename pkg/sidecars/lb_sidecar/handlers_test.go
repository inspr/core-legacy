package sidecarserv

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

	"github.com/inspr/inspr/cmd/insprd/memory/tree"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/sidecars/models"
)

func createMockEnvVars() {
	customEnvValues := "chan1;chan2;chan3;chan4;chan5"
	os.Setenv("INSPR_INPUT_CHANNELS", customEnvValues)
	os.Setenv("INSPR_OUTPUT_CHANNELS", customEnvValues)
	os.Setenv("chan1_BROKER", "invBroker2")
	os.Setenv("chan5_BROKER", "randBroker1")
	os.Setenv("INSPR_SIDECAR_INVBROKER2_WRITE_PORT", "")
	os.Setenv("INSPR_SIDECAR_RANDBROKER1_WRITE_PORT", "1107")
	os.Setenv("INSPR_LBSIDECAR_WRITE_PORT", "1127")
	os.Setenv("INSPR_LBSIDECAR_READ_PORT", "1137")
}

func deleteMockEnvVars() {
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("chan1_BROKER")
	os.Unsetenv("chan5_BROKER")
	os.Unsetenv("INSPR_SIDECAR_INVBROKER2_WRITE_PORT")
	os.Unsetenv("INSPR_SIDECAR_RANDBROKER1_WRITE_PORT")
	os.Unsetenv("INSPR_LBSIDECAR_WRITE_PORT")
	os.Unsetenv("INSPR_LBSIDECAR_READ_PORT")
}

func createMockedApp() *meta.App {
	root := meta.App{
		Meta: meta.Metadata{
			Name:   "",
			Parent: "",
			UUID:   "",
		},
		Spec: meta.AppSpec{
			Apps: map[string]*meta.App{
				"app1": {
					Meta: meta.Metadata{
						Name: "app1",
					},
					Spec: meta.AppSpec{
						Channels: map[string]*meta.Channel{
							"chan1": {
								Meta: meta.Metadata{
									Name: "chan1",
								},
								Spec: meta.ChannelSpec{
									SelectedBroker: "invBroker1",
								},
							},
							"chan4": {
								Meta: meta.Metadata{
									Name: "chan4",
								},
								Spec: meta.ChannelSpec{
									SelectedBroker: "invBroker2",
								},
							},
							"chan5": {
								Meta: meta.Metadata{
									Name: "chan5",
								},
								Spec: meta.ChannelSpec{
									SelectedBroker: "randBroker1",
								},
							},
						},
					},
				},
			},
		},
	}
	return &root
}

func createMockedServer(port, ch, msg string) *httptest.Server {
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

			if (ch != channel) || (msg != receivedData.Message) {
				rest.ERROR(w, fmt.Errorf("invalid channel or message"))
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

	treeRoot := createMockedApp()

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
			name:    "Channel has invalid broker (port envvar not found)",
			channel: "chan1",
			wantErr: true,
		},
		{
			name:    "Invalid sidecar port to send request",
			channel: "chan4",
			wantErr: true,
		},
		{
			name:    "Valid write request",
			channel: "chan5",
			msg: models.BrokerMessage{
				Message: "randomMessage",
			},
			port: "1107",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testServer *httptest.Server
			if tt.port != "" {
				testServer = createMockedServer(tt.port, tt.channel, tt.msg.Message.(string))
				testServer.Start()
				defer testServer.Close()
			}

			tree.SetMockedTree(treeRoot, nil, false, false, false)

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

	treeRoot := createMockedApp()

	rServer := httptest.NewServer(Init().readMessageHandler())
	req := http.Client{}

	tests := []struct {
		name    string
		channel string
		msg     models.BrokerMessage
		port    string
		wantErr bool
	}{
		{
			name:    "Channel not listed in 'INSPR_INPUT_CHANNELS'",
			channel: "invalidChan1",
			wantErr: true,
		},
		{
			name:    "Env var 'INSPR_SCCLIENT_READ_PORT' doesn't exist",
			channel: "chan2",
			wantErr: true,
		},
		{
			name:    "Valid read request",
			channel: "chan5",
			msg: models.BrokerMessage{
				Message: "randomMessage",
			},
			port: "1117",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				os.Setenv("INSPR_SCCLIENT_READ_PORT", "1117")
				defer os.Unsetenv("INSPR_SCCLIENT_READ_PORT")
			}
			var testServer *httptest.Server
			if tt.port != "" {
				testServer = createMockedServer(tt.port, tt.channel, tt.msg.Message.(string))
				testServer.Start()
				defer testServer.Close()
			}

			tree.SetMockedTree(treeRoot, nil, false, false, false)

			buf, _ := json.Marshal(tt.msg)
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
				fmt.Printf("Response body: %v", resp.Body)
				return
			}
		})
	}
}
