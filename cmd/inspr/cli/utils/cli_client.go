package utils

import (
	"encoding/json"
	"io"
	"os"

	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

type cliDefaults struct {
	client *client.Client
	out    io.Writer
}

var defaults cliDefaults

//GetCliClient returns the default controller client for cli.
func GetCliClient() *client.Client {
	if defaults.client == nil {
		SetDefaultClient()
	}
	return defaults.client
}

//GetCliOut returns the default output for cli.
func GetCliOut() io.Writer {
	if defaults.out == nil {
		setDefaultOut()
	}
	return defaults.out
}

//SetDefaultClient creates cli's controller client from viper's configured serverIp
func SetDefaultClient() {
	url := GetConfiguredServerIP()
	rc := request.NewClient().BaseURL(url).Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build()

	defaults = cliDefaults{
		client: client.NewControllerClient(rc),
	}
}

func setDefaultOut() {
	defaults = cliDefaults{
		out: os.Stdout,
	}
}
