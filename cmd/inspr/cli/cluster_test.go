package cli

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/inspr/inspr/cmd/sidecars"
	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/auth"
	cliutils "github.com/inspr/inspr/pkg/cmd/utils"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
					t.Errorf("NewClusterCommand() not created successfully")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClusterCommand()
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

			cmd.Execute()
			got := buf.String()

			if !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf("getBrokers() = %v, want %v", got, tt.expectedOutput)
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

			cmd.Execute()
			got := buf.String()

			if !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf("authInit() = %v, want %v", got, tt.expectedOutput)
			}

		})
	}
}

func Test_clusterConfig(t *testing.T) {
	prepareToken(t)
	defer restartScopeFlag()

	dir := os.TempDir()

	// kafka yml preparation
	kafkaFile := dir + "/kafkaConfig.yml"
	kafkaConfigBytes, _ := yaml.Marshal(sidecars.KafkaConfig{
		BootstrapServers: "mock_bootstrap",
		AutoOffsetReset:  "mock_autooffset",
		SidecarImage:     "mock_sidecarimg",
		KafkaInsprAddr:   "mock_kafkaInsprAddr",
	})
	os.WriteFile(kafkaFile, kafkaConfigBytes, 0777)
	defer os.Remove(kafkaFile)

	// invalid yml preparation
	invalidFile := dir + "/invalidConfig.yml"
	os.WriteFile(invalidFile, []byte{1}, 0777)
	defer os.Remove(invalidFile)

	// non yml file preparation
	nonYamlFile := dir + "/nonYaml.txt"
	os.WriteFile(nonYamlFile, []byte{1}, 0777)
	defer os.Remove(nonYamlFile)

	// non existant file err message
	_, err := os.Stat(dir + "nonExistantFile.yml")
	nonExistantMessage := err.Error()

	// mock error for the mockClient
	clientMockErr := errors.New("mock_error")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantMsg string
	}{
		{
			name:    "filePath_empty",
			args:    []string{"config", "kafka", ""},
			wantMsg: "arg 'filePath' is empty",
			wantErr: true,
		},
		{
			name:    "configName_empty",
			args:    []string{"config", "", invalidFile},
			wantMsg: "arg 'brokerName' is empty",
			wantErr: true,
		},
		{
			name:    "non_existant_config_file",
			args:    []string{"config", "kafka", dir + "nonExistantFile.yml"},
			wantMsg: nonExistantMessage,
			wantErr: true,
		},
		{
			name:    "non_yaml_config_file",
			args:    []string{"config", "kafka", nonYamlFile},
			wantMsg: "not a yaml file",
			wantErr: true,
		},
		{
			name:    "failed_client_create",
			args:    []string{"config", "kafka", kafkaFile},
			wantMsg: clientMockErr.Error(),
			wantErr: true,
		},
		{
			name:    "working",
			args:    []string{"config", "kafka", kafkaFile},
			wantMsg: "successfully installed broker on insprd\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewClusterCommand()
			buf := bytes.NewBufferString("")

			cliutils.SetOutput(buf)
			if !tt.wantErr {
				cliutils.SetMockedClient(nil)
			} else {
				cliutils.SetMockedClient(clientMockErr)
			}

			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			// getting the output generated by the cli
			out := buf.String()

			if out != tt.wantMsg {
				t.Errorf("cluster config msg error, got '%v' want '%v'", out, tt.wantMsg)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("cluster config error, got '%v' want '%v'", err, tt.wantErr)
			}
		})
	}
}
