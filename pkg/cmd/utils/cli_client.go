package utils

import (
	"io"
	"os"
	"path/filepath"

	"inspr.dev/inspr/pkg/cmd"
	"inspr.dev/inspr/pkg/controller"
	"inspr.dev/inspr/pkg/controller/client"
	"inspr.dev/inspr/pkg/controller/mocks"
)

type cliGlobalStructure struct {
	client controller.Interface
	out    io.Writer
}

var defaults cliGlobalStructure

//GetCliClient returns the default controller client for cli.
func GetCliClient() controller.Interface {
	if defaults.client == nil {
		SetGlobalClient()
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

//setGlobalClient creates cli's controller client from viper's configured serverIp
func SetGlobalClient() {
	url := GetConfiguredServerIP()
	host := GetConfiguredHost()

	if cmd.InsprOptions.Host != "" {
		host = cmd.InsprOptions.Host
	}

	SetClient(url, host)
}

func setGlobalOutput() {
	defaults.out = os.Stdout
}

// SetOutput sets the default output of CLI
func SetOutput(out io.Writer) {
	defaults.out = out
}

// SetClient sets the default server IP of CLI
func SetClient(url string, host string) {
	if cmd.InsprOptions.Token == "" {
		dir, _ := os.UserHomeDir()
		cmd.InsprOptions.Token = filepath.Join(dir, ".inspr/token")
	}

	config := client.ControllerConfig{
		Auth: Authenticator{
			cmd.InsprOptions.Token,
		},
		URL:  url,
		Host: host,
	}

	defaults.client = client.NewControllerClient(config)
}

//SetMockedClient configures singleton's client as a mocked client given a error
func SetMockedClient(err error) {
	defaults.client = mocks.NewClientMock(err)
}
