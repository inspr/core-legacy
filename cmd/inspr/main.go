package main

import (
	"os"

	"github.com/inspr/inspr/cmd/inspr/cli"

	"github.com/inspr/inspr/pkg/meta"
)

var version string

func init() {
	if version == "" {
		version = "not given"
	}
}

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

	cli.GetFactory().Subscribe(meta.Component{
		APIVersion: "v1",
		Kind:       "alias",
	}, cli.NewApplyAlias())

	cli.NewInsprCommand(os.Stdout, os.Stderr, version).Execute()
}
