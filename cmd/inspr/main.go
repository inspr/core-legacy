package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	cli "gitlab.inspr.dev/inspr/core/cmd/inspr/cli"
	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
	"gopkg.in/yaml.v2"
)

func main() {
	rc := request.NewClient().
		BaseURL("http://127.0.0.1:8080").
		Encoder(json.Marshal).
		Decoder(request.JSONDecoderGenerator).
		Build()
	client := client.NewControllerClient(rc)

	cli.GetFactory().Subscribe(meta.Component{
		APIVersion: "v1",
		Kind:       "channel",
	},
		func(b []byte, out io.Writer) error {
			ch := meta.Channel{}

			yaml.Unmarshal(b, &ch)
			fmt.Println(ch)

			return nil
		})

	cli.GetFactory().Subscribe(meta.Component{
		APIVersion: "v1",
		Kind:       "channeltype",
	}, cli.NewApplyChannelType(client.ChannelTypes()))
	cli.NewInsprCommand(os.Stdout, os.Stderr).Execute()
}
