package utils

import (
	"encoding/json"
	"io"
	"os"

	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

type cliGlobalStructure struct {
	client *client.Client
	out    io.Writer
}

var defaults cliGlobalStructure

//GetCliClient returns the default controller client for cli.
func GetCliClient() *client.Client {
	if defaults.client == nil {
		setGlobalClient()
	}
	return defaults.client
}

//GetCliOut returns the default output for cli.
func GetCliOutput() io.Writer {
	if defaults.out == nil {
		setGlobalOutput()
	}
	return defaults.out
}

//SetDefaultClient creates cli's controller client from viper's configured serverIp
func setGlobalClient() {
	url := GetConfiguredServerIP()
	rc := request.NewClient().BaseURL(url).Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build()

	defaults = cliGlobalStructure{
		client: client.NewControllerClient(rc),
	}
}

func setGlobalOutput() {
	defaults = cliGlobalStructure{
		out: os.Stdout,
	}
}
