package main

import (
	"fmt"
	"io"
	"os"

	cli "gitlab.inspr.dev/inspr/core/cmd/inspr/cli"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gopkg.in/yaml.v2"
)

func main() {
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
	cli.NewInsprCommand(os.Stdout, os.Stderr).Execute()
}
