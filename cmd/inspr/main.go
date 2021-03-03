package main

import (
	"os"

	cli "gitlab.inspr.dev/inspr/core/cmd/inspr/cli"
	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func main() {
	client := cliutils.GetCliClient()

	cli.GetFactory().Subscribe(meta.Component{
		APIVersion: "v1",
		Kind:       "channel",
	}, cli.NewApplyChannel(client.Channels()))

	cli.GetFactory().Subscribe(meta.Component{
		APIVersion: "v1",
		Kind:       "channeltype",
	}, cli.NewApplyChannelType(client.ChannelTypes()))

	cli.GetFactory().Subscribe(meta.Component{
		APIVersion: "v1",
		Kind:       "dapp",
	}, cli.NewApplyApp(client.Apps()))

	cli.NewInsprCommand(os.Stdout, os.Stderr).Execute()
}
