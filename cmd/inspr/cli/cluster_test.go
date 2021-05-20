package cli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/auth"
	cliutils "github.com/inspr/inspr/pkg/cmd/utils"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/spf13/cobra"
)

func TestNewClusterCommand(t *testing.T) {
	tests := []struct {
		name          string
		checkFunction func(t *testing.T, got *cobra.Command)
	}{
		{
			name: "It should create a new delete command",
			checkFunction: func(t *testing.T, got *cobra.Command) {
				if got == nil {
					t.Errorf("NewDeleteCmd() not created successfully")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDeleteCmd()
			if tt.checkFunction != nil {
				tt.checkFunction(t, got)
			}
		})
	}
}

func Test_getBrokers(t *testing.T) {
	prepareToken(t)
	defer restartScopeFlag()

	tests := []struct {
		name           string
		flagsAndArgs   []string
		expectedOutput string
		handlerFunc    func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:         "Should get brokers",
			flagsAndArgs: []string{"b"},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				rest.JSON(w, 200, &models.BrokersDI{
					Installed: []string{"mock_broker"},
					Default:   "mock_broker",
				})
			},
			expectedOutput: "DEFAULT:\nmock_broker\nAVAILABLE:\nmock_broker\n",
		},
		{
			name:         "Should return error",
			flagsAndArgs: []string{"b"},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.NewError().Message("error").Build())
			},
			expectedOutput: "unexpected inspr error, the message is: error\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewClusterCommand()

			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)

			cmd.SetArgs(tt.flagsAndArgs)

			server := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			cliutils.SetClient(server.URL)

			defer server.Close()

			bufResp := bytes.NewBufferString("")
			fmt.Fprint(bufResp, tt.expectedOutput)

			outResp, _ := ioutil.ReadAll(bufResp)

			cmd.Execute()
			got, _ := ioutil.ReadAll(buf)

			if !reflect.DeepEqual(got, outResp) {
				t.Errorf("getBrokers() = %v, want %v", string(got), string(outResp))
			}
		})
	}
}

func Test_authInit(t *testing.T) {
	prepareToken(t)
	defer restartScopeFlag()

	tests := []struct {
		name           string
		flagsAndArgs   []string
		expectedOutput string
		handlerFunc    func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:         "Should get token",
			flagsAndArgs: []string{"init", "pwd"},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				rest.JSON(w, 200, &auth.JwtDO{
					Token: []byte("mock_token"),
				})
			},
			expectedOutput: "This is a root token for authentication within your insprd. This will not be generated again. Save it wisely.\nmock_token\n",
		},
		{
			name:         "Should return error",
			flagsAndArgs: []string{"init", "pwd"},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.NewError().Message("error").Build())
			},
			expectedOutput: "unexpected inspr error, the message is: error\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewClusterCommand()

			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)

			cmd.SetArgs(tt.flagsAndArgs)

			server := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			cliutils.SetClient(server.URL)

			defer server.Close()

			bufResp := bytes.NewBufferString("")
			fmt.Fprint(bufResp, tt.expectedOutput)

			outResp, _ := ioutil.ReadAll(bufResp)

			cmd.Execute()
			got, _ := ioutil.ReadAll(buf)

			if !reflect.DeepEqual(got, outResp) {
				t.Errorf("authInit() = %v, want %v", string(got), string(outResp))
			}

		})
	}
}
