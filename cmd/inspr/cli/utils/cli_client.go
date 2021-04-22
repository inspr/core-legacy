package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/inspr/inspr/pkg/controller"
	"github.com/inspr/inspr/pkg/controller/client"
	"github.com/inspr/inspr/pkg/controller/mocks"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest/request"
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
	rc := request.NewJSONClient(url)

	defaults.client = client.NewControllerClient(rc)
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
	rc := request.NewClient().BaseURL(url).Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build()
	defaults.client = client.NewControllerClient(rc)
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
			fmt.Fprintf(w, "Did you login ?\n")
		case ierrors.Forbidden:
			fmt.Fprintf(w, "Forbidden operation, please check for the scope.\n")
		default:
			fmt.Fprintf(w, "unexpected inspr error, the message is: %v\n", err.Error())
		}
	} else {
		fmt.Fprintf(w, "Non inspr error, the message is: %v\n", err.Error())
	}
}
