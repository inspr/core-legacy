package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/inspr/inspr/pkg/cmd"
	"github.com/inspr/inspr/pkg/controller"
	"github.com/inspr/inspr/pkg/controller/client"
	"github.com/inspr/inspr/pkg/controller/mocks"
	"github.com/inspr/inspr/pkg/ierrors"
)

type cliGlobalStructure struct {
	client controller.Interface
	out    io.Writer
}

var defaults cliGlobalStructure

//GetCliClient returns the default controller client for cli.
func GetCliClient() controller.Interface {
	if defaults.client == nil {
		setGlobalClient()
	}
	return defaults.client
}

//GetCliOutput returns the default output for cli.
func GetCliOutput() io.Writer {
	if defaults.out == nil {
		setGlobalOutput()
	}
	return defaults.out
}

//SetDefaultClient creates cli's controller client from viper's configured serverIp
func setGlobalClient() {
	url := GetConfiguredServerIP()
	SetClient(url)
}

func setGlobalOutput() {
	defaults.out = os.Stdout
}

// SetOutput sets the default output of CLI
func SetOutput(out io.Writer) {
	defaults.out = out
}

// SetClient sets the default server IP of CLI
func SetClient(url string) {
	if cmd.InsprOptions.Token == "" {
		dir, _ := os.UserHomeDir()
		cmd.InsprOptions.Token = filepath.Join(dir, ".inspr/token")
	}

	config := client.ControllerConfig{
		Auth: Authenticator{
			cmd.InsprOptions.Token,
		},
		URL: url,
	}

	defaults.client = client.NewControllerClient(config)
}

//SetMockedClient configures singleton's client as a mocked client given a error
func SetMockedClient(err error) {
	defaults.client = mocks.NewClientMock(err)
}

// RequestErrorMessage prints an error to the user based on the error given, in
// actuality it converts the error to an insprErr and then process what type
// of return the apply request returned.
func RequestErrorMessage(err error, w io.Writer) {
	ierr, ok := err.(*ierrors.InsprError)
	if ok {
		switch ierr.Code {
		case ierrors.Unauthorized:
			fmt.Fprintf(w, "we couldn't authenticate with the cluster. Is your token configured correctly?\n")
		case ierrors.Forbidden:
			fmt.Fprintf(w, "forbidden operation, please check for the scope.\n")
		default:
			fmt.Fprintf(w, "unexpected inspr error: %v\n", err.Error())
		}
	} else {
		fmt.Fprintf(w, "non inspr error: %v\n", err.Error())
	}
}
