package main

import (
	"os"

	cli "gitlab.inspr.dev/inspr/core/cmd/inspr/cli"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

var version string

func main() {
	cli.GetFactory().Subscribe(meta.Component{
		APIVersion: "v1",
		Kind:       "channel",
	}, cli.NewApplyChannel())

	cli.GetFactory().Subscribe(meta.Component{
		APIVersion: "v1",
		Kind:       "channeltype",
	}, cli.NewApplyChannelType())

	cli.GetFactory().Subscribe(meta.Component{
		APIVersion: "v1",
		Kind:       "dapp",
	}, cli.NewApplyApp())

	cli.NewInsprCommand(os.Stdout, os.Stderr, version).Execute()
}
