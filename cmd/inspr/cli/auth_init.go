package cli

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/inspr/inspr/pkg/cmd"
	"github.com/inspr/inspr/pkg/cmd/utils"
	"github.com/inspr/inspr/pkg/ierrors"
)

type clusterInitConfigDT struct {
	providerKind string
	providerYAML string
	providerJSON string
}

var clusterInitConfig clusterInitConfigDT

var authCommand = cmd.NewCmd("cluster").AddSubCommand(
	cmd.NewCmd("init").
		WithFlags(
			[]*cmd.Flag{
				{
					Name:     "yaml",
					Value:    &clusterInitConfig.providerYAML,
					DefValue: "",
				},
				{
					Name:     "json",
					Value:    &clusterInitConfig.providerJSON,
					DefValue: "",
				},
			},
		).
		ExactArgs(1,
			func(c context.Context, args []string) error {
				var err error
				provBytes := []byte{}
				if clusterInitConfig.providerJSON != "" {
					provBytes, err = ioutil.ReadFile(clusterInitConfig.providerJSON)
					if err != nil {
						return err
					}
				} else if clusterInitConfig.providerYAML != "" {

					provBytes, err = ioutil.ReadFile(clusterInitConfig.providerYAML)
					if err != nil {
						return err
					}
				} else {
					return ierrors.NewError().Message("you need to set a provider via either the yaml or json flags.").Build()
				}
				token, err := utils.GetCliClient().Authorization().Init(c, args[0], provBytes)
				if err != nil {
					return err
				}
				output := utils.GetCliOutput()

				fmt.Fprintln(output, "This is a root token for authentication within your insprd. This will not be generated again. Save it wisely.")
				fmt.Fprintf(output, "%s\n", token)
				return nil
			},
		),
).WithCommonFlags().Super()
