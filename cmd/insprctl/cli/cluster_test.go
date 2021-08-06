package cli

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"inspr.dev/inspr/pkg/auth"
	cliutils "inspr.dev/inspr/pkg/cmd/utils"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/rest"
	"inspr.dev/inspr/pkg/rest/request"
)

func TestNewClusterCommand(t *testing.T) {
	tests := []struct {
		name          string
		checkFunction func(t *testing.T, got *cobra.Command)
	}{
		{
			name: "It should create a new cluster command",
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
			name:         "Should_return_body_error",
			flagsAndArgs: []string{"init", "pwd"},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.New("error"))
			},
			expectedOutput: ierrors.FormatError(ierrors.New("error")),
		},
		{
			name:         "Should_return_default_request_error",
			flagsAndArgs: []string{"init", "pwd"},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, nil)
			},
			expectedOutput: ierrors.FormatError(request.DefaultErr),
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
